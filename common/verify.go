package common

import (
	"bytes"
	"crypto/ed25519"
	"crypto/sha256"
	"encoding/base64"
	"encoding/hex"
	"errors"
	"fmt"
	"regexp"
	"strings"

	"github.com/Conflux-Chain/go-conflux-sdk/types/cfxaddress"
	"github.com/martinboehm/btcd/txscript"
	"github.com/martinboehm/btcutil"
	"github.com/martinboehm/btcutil/base58"
	"github.com/martinboehm/btcutil/bech32"
	"github.com/martinboehm/btcutil/chaincfg"
	"golang.org/x/crypto/blake2b"
	"golang.org/x/crypto/ripemd160"

	secp_ecdsa "github.com/decred/dcrd/dcrec/secp256k1/v4/ecdsa"
	"github.com/okx/go-wallet-sdk/coins/cosmos"
	"github.com/okx/go-wallet-sdk/coins/stacks"
	tonWallet "github.com/xssnick/tonutils-go/ton/wallet"
	"golang.org/x/crypto/sha3"
)

var (
	ErrInvalidAddr = errors.New("invalid address")
	ErrInvalidSign = errors.New("can't verify signature")
)

func VerifyBETH(addr, msg, sign string) error {
	coin := "BETH"
	msgHeader, exist := PorCoinMessageSignatureHeaderMap[coin]
	if !exist {
		return fmt.Errorf("invalid coin type %s", coin)
	}
	hash := HashEvmCoinTypeMsg(msgHeader, msg)
	var p [48]byte
	var h [32]byte
	var s [96]byte
	copy(p[:], MustDecode(addr))
	copy(s[:], MustDecode(sign))
	copy(h[:], hash[:])
	var pub Pubkey
	if err := pub.Deserialize(&p); err != nil {
		return fmt.Errorf("unexpected failure, failed to deserialize pubkey (%x): %v", p[:], err)
	}
	var sig Signature
	if err := sig.Deserialize(&s); err != nil {
		return fmt.Errorf("unexpected failure, failed to deserialize signature (%x): %v", s[:], err)
	}
	res := Verify(&pub, h[:], &sig)
	if !res {
		return errors.New("unexpected failure, failed to verify signature")
	}
	return nil
}

func VerifyTRX(addr, msg, sign string) error {
	hashFuncs := []func(string) []byte{HashTrxMsg, HashTrxMsgV2}

	for _, hashFunc := range hashFuncs {
		if verifyTRX(addr, msg, sign, hashFunc) == nil {
			return nil
		}
	}

	return ErrInvalidSign
}

func verifyTRX(addr, msg, sign string, hashFunc func(string) []byte) error {
	hash := hashFunc(msg)
	s := MustDecode(sign)
	pub, err := sigToPub(hash, s)
	if err != nil {
		return fmt.Errorf("failed to recover public key from TRX signature, error:%v", err)
	}
	pubKey := pub.SerializeUncompressed()
	h := sha3.NewLegacyKeccak256()
	h.Write(pubKey[1:])
	newHash := h.Sum(nil)[12:]
	newAddr := base58.CheckEncode(newHash, GetNetWork(), base58.Sha256D)
	if addr != newAddr {
		return fmt.Errorf("TRX address mismatch, expected:%s, recovered:%s", addr, newAddr)
	}
	return nil
}

func UtxoCoinSigToPubKey(coin, msg, sign string) ([]byte, error) {
	msgHeader, exist := PorCoinMessageSignatureHeaderMap[coin]
	if !exist {
		return nil, fmt.Errorf("invalid coin type %s", coin)
	}
	hash := HashUtxoCoinTypeMsg(msgHeader, msg)
	b, err := base64.StdEncoding.DecodeString(sign)
	if err != nil {
		return nil, fmt.Errorf("failed to decode UTXO signature, coin:%s, error:%v", coin, err)
	}
	pub, ok, err := secp_ecdsa.RecoverCompact(b, hash)
	if err != nil || !ok || pub == nil {
		return nil, fmt.Errorf("failed to recover UTXO public key from signature, coin:%s, error:%v", coin, err)
	}

	return pub.SerializeCompressed(), nil
}

func VerifyUtxoCoin(coin, addr, msg, sign1, sign2, script string) error {
	var pub1, pub2 []byte
	var err error
	// recover pub1 and pub2 from sign1 and sign2
	if sign1 != "" {
		pub1, err = UtxoCoinSigToPubKey(coin, msg, sign1)
		if err != nil {
			return err
		}
	}
	if sign2 != "" {
		pub2, err = UtxoCoinSigToPubKey(coin, msg, sign2)
		if err != nil {
			return err
		}
	}

	return VerifyUtxoCoinSig(coin, addr, script, pub1, pub2)
}

func VerifyUtxoCoinSig(coin, addr, script string, pub1, pub2 []byte) error {
	mainNetParams := &chaincfg.Params{}
	coinAddressType := PorCoinAddressTypeMap[coin]
	// get main net params
	switch coinAddressType {
	case "BTC":
		mainNetParams = GetBTCMainNetParams()
	case "BCH":
		mainNetParams = GetBTCMainNetParams()
		// convert cash address to legacy address
		if IsCashAddress(addr) {
			legacyAddr, err := ConvertCashAddressToLegacy(addr)
			if err != nil {
				return fmt.Errorf("convertCashAddressToLegacy failed, invalid cash address: %s, error: %v", addr, err)
			}
			addr = legacyAddr
		}
	case "LTC":
		mainNetParams = GetLTCMainNetParams()
	case "DOGE":
		mainNetParams = GetDOGEMainNetParams()
	case "DASH":
		mainNetParams = GetDASHMainNetParams()
	case "BTG":
		mainNetParams = GetBTGMainNetParams()
	case "DGB":
		mainNetParams = GetDGBMainNetParams()
	case "QTUM":
		mainNetParams = GetQTUMMainNetParams()
	case "RVN":
		mainNetParams = GetRVNMainNetParams()
	case "ZEC":
		mainNetParams = GetZECMainNetParams()
	default:
		mainNetParams = GetBTCMainNetParams()
	}
	if _, err := btcutil.DecodeAddress(addr, mainNetParams); err != nil {
		return fmt.Errorf("invalid UTXO address format, coin:%s, addr:%s, error:%v", coin, addr, err)
	}
	addrType := GuessUtxoCoinAddressType(addr)
	switch addrType {
	case "P2PKH":
		addrPub, err := btcutil.NewAddressPubKey(pub1, mainNetParams)
		if err != nil || addrPub.EncodeAddress() != addr {
			return fmt.Errorf("address not match,coin: %s, addr: %s, recoverAddr: %s", coin, addr, addrPub.EncodeAddress())
		}
	case "P2SH":
		if script == "" {
			return fmt.Errorf("P2SH address requires script, but script is empty, coin:%s, addr:%s", coin, addr)
		}
		addrPub, err := btcutil.NewAddressScriptHash(MustDecode(script), mainNetParams)
		if err != nil {
			return fmt.Errorf("get NewAddressScriptHash failed, coin:%s, addr:%s, error:%v", coin, addr, err)
		}
		if addrPub.EncodeAddress() != addr {
			return fmt.Errorf("address not match, coin:%s, addr:%s, recoverAddr:%s", coin, addr, addrPub.EncodeAddress())
		}
		addrPub1, err := btcutil.NewAddressPubKey(pub1, mainNetParams)
		if err != nil {
			return fmt.Errorf("get pub1 NewAddressPubKey failed, coin:%s, addr:%s, error: %v", coin, addr, err)
		}
		addr1 := addrPub1.EncodeAddress()
		addrPub2, err := btcutil.NewAddressPubKey(pub2, mainNetParams)
		if err != nil {
			return fmt.Errorf("get pub2 NewAddressPubKey failed, coin:%s, addr:%s, error:%v", coin, addr, err)
		}
		addr2 := addrPub2.EncodeAddress()
		typ, pubs, _, err := txscript.ExtractPkScriptAddrs(MustDecode(script), mainNetParams)
		if typ != txscript.MultiSigTy {
			return fmt.Errorf("script type not match, coin:%s, addr:%s, srcType:%d, type:%d", coin, addr, txscript.MultiSigTy, typ)
		}
		if err != nil {
			return fmt.Errorf("script ExtractPkScriptAddrs failed, coin:%s, addr:%s, error: %v", coin, addr, err)
		}
		m := map[string]struct{}{addr1: {}, addr: {}, addr2: {}}
		for _, v := range pubs {
			delete(m, v.EncodeAddress())
		}
		if len(m) > 1 {
			return fmt.Errorf("script address not match the pubs, coin:%s, addr:%s", coin, addr)
		}
	case "P2WSH":
		pkScript := MustDecode(script)
		h := sha256.New()
		h.Write(pkScript)
		witnessProg := h.Sum(nil)
		addressWitnessScriptHash, err := btcutil.NewAddressWitnessScriptHash(witnessProg, mainNetParams)
		if err != nil {
			return fmt.Errorf("get NewAddressWitnessScriptHash failed, coin:%s, addr:%s, error:%v", coin, addr, err)
		}
		if addressWitnessScriptHash.EncodeAddress() != addr {
			return fmt.Errorf("address not match,coin: %s, addr: %s, recoverAddr: %s", coin, addr, addressWitnessScriptHash.EncodeAddress())
		}
		addrPub1, err := btcutil.NewAddressPubKey(pub1, mainNetParams)
		if err != nil {
			return fmt.Errorf("get pub1 NewAddressPubKey failed, coin:%s, addr:%s, error: %v", coin, addr, err)
		}
		addr1 := addrPub1.EncodeAddress()
		addrPub2, err := btcutil.NewAddressPubKey(pub2, mainNetParams)
		if err != nil {
			return fmt.Errorf("get pub2 NewAddressPubKey failed, coin:%s, addr:%s, error: %v", coin, addr, err)
		}
		addr2 := addrPub2.EncodeAddress()
		typ, pubs, _, err := txscript.ExtractPkScriptAddrs(MustDecode(script), mainNetParams)
		if typ != txscript.MultiSigTy {
			return fmt.Errorf("script type not match, coin:%s, addr:%s, srcType:%d, type:%d", coin, addr, txscript.MultiSigTy, typ)
		}
		if err != nil {
			return fmt.Errorf("script ExtractPkScriptAddrs failed, coin:%s, addr:%s, error: %v", coin, addr, err)
		}
		m := map[string]struct{}{addr1: {}, addr: {}, addr2: {}}
		for _, v := range pubs {
			delete(m, v.EncodeAddress())
		}
		if len(m) > 1 {
			return fmt.Errorf("script address not match the pubs, coin:%s, addr:%s", coin, addr)
		}
	}
	return nil
}

func VerifyEvmCoin(coin, addr, msg, sign string) error {
	msgHeader, exist := PorCoinMessageSignatureHeaderMap[coin]
	if !exist {
		return fmt.Errorf("invalid coin type %s, addr:%s", coin, addr)
	}
	hash := HashEvmCoinTypeMsg(msgHeader, msg)
	s := MustDecode(sign)
	pub, err := sigToPub(hash, s)
	if err != nil {
		return ErrInvalidAddr
	}

	pubToEcdsa := pub.ToECDSA()
	recoverAddr := PubkeyToAddress(*pubToEcdsa).String()

	addrType, exist := PorCoinAddressTypeMap[coin]
	if !exist {
		return fmt.Errorf("invalid coin type %s, addr:%s", coin, addr)
	}
	switch addrType {
	case "FIL":
		// convert eth address to fil address
		filAddress, err := ConvertEthAddressToFilecoinAddress(PubkeyToAddress(*pubToEcdsa).Bytes())
		if err != nil {
			return fmt.Errorf("convert eth address to fil address failed, coin:%s, addr:%s, error:%v", coin, addr, err)
		}
		recoverAddr = filAddress.String()
	case "LAT":
		ethAddress := PubkeyToAddress(*pubToEcdsa).String()
		recoverAddr, _ = ConvertETHToLATAddress(ethAddress)

	case "ETH":
		if !VerifySignAddr(HexToAddress(addr), hash, s) {
			// 获取恢复出来的地址用于错误信息
			recoveredAddr := PubkeyToAddress(*pubToEcdsa).String()
			return fmt.Errorf("ETH address verification failed, coin:%s, expected:%s, recovered:%s", coin, addr, recoveredAddr)
		}
	}

	if !strings.EqualFold(addr, recoverAddr) {
		return fmt.Errorf("recovery address not match, coin:%s, recoverAddr:%s, addr:%s", coin, recoverAddr, addr)
	}

	return nil
}

func VerifyEd25519Coin(coin, addr, msg, sign, pubkey string) error {
	msgHeader, exist := PorCoinMessageSignatureHeaderMap[coin]
	if !exist {
		return fmt.Errorf("invalid coin type %s, addr:%s", coin, addr)
	}
	hash := HashEd25519Msg(msgHeader, msg)
	res, _ := Decode(sign)
	pubkeyBytes, _ := Decode(pubkey)
	if ok := ed25519.Verify(pubkeyBytes, hash, res); !ok {
		return fmt.Errorf("ED25519 signature verification failed, coin:%s, addr:%s", coin, addr)
	}

	addrType, exist := PorCoinAddressTypeMap[coin]
	if !exist {
		return fmt.Errorf("invalid coin type %s, addr:%s", coin, addr)
	}
	var recoverAddrs []string
	switch addrType {
	case "SOL":
		out := [32]byte{}
		byteCount := len(pubkeyBytes)
		if byteCount == 0 {
			return fmt.Errorf("empty public key for SOL address generation, coin:%s, addr:%s", coin, addr)
		}
		max := 32
		if byteCount < max {
			max = byteCount
		}
		copy(out[:], pubkeyBytes[0:max])
		recoverAddrs = append(recoverAddrs, base58.Encode(out[:]))
	case "APTOS":
		publicKey := append(pubkeyBytes, 0x0)
		rAddr := "0x" + hex.EncodeToString(Sha256Hash(publicKey))
		// Short address type: if address starts with 0x0, replace.
		re, _ := regexp.Compile("^0x0*")
		recoverAddrs = append(recoverAddrs, re.ReplaceAllString(rAddr, "0x"))
		recoverAddrs = append(recoverAddrs, rAddr)
	case "SUI":
		k := make([]byte, 33)
		copy(k[1:], pubkeyBytes)
		publicKeyHash, err := blake2b.New256(nil)
		if err != nil {
			return fmt.Errorf("invalid publicKey, coin:%s, recoverAddrs:%v, addr:%s", coin, recoverAddrs, addr)
		}
		publicKeyHash.Write(k)
		h := publicKeyHash.Sum(nil)
		address := "0x" + hex.EncodeToString(h)[0:64]
		recoverAddrs = append(recoverAddrs, address)
	case "TON":
		walletV3, err := tonWallet.AddressFromPubKey(pubkeyBytes, tonWallet.V3, tonWallet.DefaultSubwallet)
		if err != nil {
			return fmt.Errorf("%s, coin: %s, addr: %s, error: %v", ErrInvalidSign, coin, addr, err)
		}
		recoverAddrs = append(recoverAddrs, walletV3.String())
		recoverAddrs = append(recoverAddrs, walletV3.Bounce(false).String())

		walletHighload, err := tonWallet.AddressFromPubKey(pubkeyBytes, tonWallet.ConfigHighloadV3{MessageTTL: 60 * 60 * 12}, 4269)
		if err != nil {
			return fmt.Errorf("%s, coin: %s, addr: %s, error: %v", ErrInvalidSign, coin, addr, err)
		}
		recoverAddrs = append(recoverAddrs, walletHighload.String())
		recoverAddrs = append(recoverAddrs, walletHighload.Bounce(false).String())
	case "DOT":
		rAddr, err := GetDotAddressFromPublicKey(pubkey)
		if err != nil {
			return fmt.Errorf("%s, coin: %s, addr: %s, error: %v", ErrInvalidSign, coin, addr, err)
		}
		recoverAddrs = append(recoverAddrs, rAddr)
	case "XLM", "PI", "STELLAR":
		// XLM addresses are base32 encoded and start with 'G'
		// The address is derived from the public key using SHA256 and then RIPEMD160
		// followed by base32 encoding with checksum
		rAddr, err := GetXlmAddressFromPublicKey(pubkey)
		if err != nil {
			return fmt.Errorf("%s, coin: %s, addr: %s, error: %v", ErrInvalidSign, coin, addr, err)
		}
		recoverAddrs = append(recoverAddrs, rAddr)
	case "ADA":
		// ADA addresses use ed25519 + Blake2b-224 + Bech32 encoding
		rAddr, err := GetAdaAddressFromPublicKey(pubkey)
		if err != nil {
			return fmt.Errorf("%s, coin: %s, addr: %s, error: %v", ErrInvalidSign, coin, addr, err)
		}
		recoverAddrs = append(recoverAddrs, rAddr)
	case "NEAR", "HBAR", "SC", "IOTA":
		return nil
	default:
		return nil
	}

	for _, recoverAddr := range recoverAddrs {
		if strings.EqualFold(recoverAddr, addr) {
			return nil
		}
	}

	return fmt.Errorf("recovery address not match, coin:%s, recoverAddrs:%v, addr:%s", coin, recoverAddrs, addr)
}

func VerifyEcdsaCoin(coin, addr, msg, sign string) error {
	msgHeader, exist := PorCoinMessageSignatureHeaderMap[coin]
	if !exist {
		return fmt.Errorf("invalid coin type %s, addr:%s", coin, addr)
	}
	hash := HashEcdsaMsg(msgHeader, msg)
	s := MustDecode(sign)
	pub, err := sigToPub(hash, s)
	if err != nil {
		return fmt.Errorf("failed to recover public key from signature, coin:%s, addr:%s, error:%v", coin, addr, err)
	}
	pubKey := pub.SerializeUncompressed()
	pubKeyCompressed := pub.SerializeCompressed()

	var recoverAddr string
	addrType, exist := PorCoinAddressTypeMap[coin]
	if !exist {
		return fmt.Errorf("invalid coin type %s, addr:%s", coin, addr)
	}
	switch addrType {
	case "FIL":
		pubKeyHash := hash_cal(pubKey, payloadHashConfig)
		explen := 1 + len(pubKeyHash)
		buf := make([]byte, explen)
		var protocol byte = 1
		buf[0] = protocol
		copy(buf[1:], pubKeyHash)
		cksm := hash_cal(buf, checksumHashConfig)
		recoverAddr = "f" + fmt.Sprintf("%d", protocol) + AddressEncoding.WithPadding(-1).EncodeToString(append(pubKeyHash, cksm[:]...))
	case "CFX":
		pubToEcdsa, _ := UnmarshalPubkey(pubKey)
		ethAddr := PubkeyToAddress(*pubToEcdsa).String()
		cfxOldAddr := "0x1" + ethAddr[3:]
		cfxAddr, err := cfxaddress.New(cfxOldAddr, 1029)
		if err != nil {
			return ErrInvalidSign
		}
		recoverAddr = cfxAddr.String()
	case "ELF":
		firstBytes := sha256.Sum256(pubKey)
		secondBytes := sha256.Sum256(firstBytes[:])
		recoverAddr = encodeCheck(secondBytes[:])
	case "LUNC":
		sha := sha256.Sum256(pub.SerializeCompressed())
		hasherRIPEMD160 := ripemd160.New()
		hasherRIPEMD160.Write(sha[:])
		recoverAddr, _ = bech32.EncodeFromBase256("terra", hasherRIPEMD160.Sum(nil))
	case "ETH":
		// OKT cosmos address type (start with 'ex')
		if strings.HasPrefix(addr, "ex") {
			hash := sha3.NewLegacyKeccak256()
			hash.Write(pubKey[1:])
			addressByte := hash.Sum(nil)
			recoverAddr, _ = bech32.EncodeFromBase256("ex", addressByte[12:])
		} else {
			pubToEcdsa := pub.ToECDSA()
			recoverAddr = PubkeyToAddress(*pubToEcdsa).String()
		}
	case "FLOW":
		pubToEcdsa := pub.ToECDSA()
		recoverAddr = PubkeyToAddress(*pubToEcdsa).String()
	case "ICX":
		recoverAddr, _ = GenerateICXAddress(pubKey)
	case "STX":
		recoverAddr, _ = stacks.GetAddressFromPublicKey(hex.EncodeToString(pubKey))
	case "TIA":
		recoverAddr, _ = cosmos.GetAddressByPublicKey(hex.EncodeToString(pubKeyCompressed), "celestia")
	case "ATOM":
		recoverAddr, _ = cosmos.GetAddressByPublicKey(hex.EncodeToString(pubKeyCompressed), "cosmos")
	case "CRO":
		recoverAddr, _ = cosmos.GetAddressByPublicKey(hex.EncodeToString(pubKeyCompressed), "cro")
	case "DORA":
		recoverAddr, _ = cosmos.GetAddressByPublicKey(hex.EncodeToString(pubKeyCompressed), "dora")
	case "DYDX":
		recoverAddr, _ = cosmos.GetAddressByPublicKey(hex.EncodeToString(pubKeyCompressed), "dydx")
	case "TERRA":
		recoverAddr, _ = cosmos.GetAddressByPublicKey(hex.EncodeToString(pubKeyCompressed), "terra")
	case "INJ":
		hash := sha3.NewLegacyKeccak256()
		hash.Write(pubKey[1:])
		addressByte := hash.Sum(nil)
		recoverAddr, _ = bech32.EncodeFromBase256("inj", addressByte[12:])
	case "ONE":
		recoverAddr, _ = GenerateONEAddress(pubKey)
	}

	if !strings.EqualFold(recoverAddr, addr) {
		return fmt.Errorf("recovery address not match, coin:%s, recoverAddr:%s, addr:%s", coin, recoverAddr, addr)
	}

	return nil
}

func VerifyStarkCoin(coin, addr, msg, sign, publicKey string) error {
	if VerifyStarknetEIP712(addr, msg, publicKey, sign) {
		return nil
	}

	return fmt.Errorf("recovery address not match, coin:%s, addr:%s", coin, addr)
}

func VerifyEcdsaCoinWithPub(msg, sign, publicKey string) error {
	hash := HashEcdsaMsg(OKXMessageSignatureHeader, msg)
	s := MustDecode(sign)
	pub, err := sigToPub(hash, s)
	if err != nil {
		return fmt.Errorf("failed to recover public key from signature, error:%v", err)
	}

	// Convert the recovered public key to uncompressed format for comparison
	recoveredPubKey := pub.SerializeCompressed()

	// Decode the provided public key
	providedPubKey, err := Decode(publicKey)
	if err != nil {
		return fmt.Errorf("invalid public key format: %v", err)
	}

	if len(providedPubKey) == 64 {
		x := providedPubKey[:32]
		y := providedPubKey[32:]
		compressed := make([]byte, 33)
		copy(compressed[1:], x)

		if y[31]&1 == 0 {
			compressed[0] = 0x02
		} else {
			compressed[0] = 0x03
		}

		providedPubKey = compressed
	}

	// Compare the public keys directly
	if !bytes.Equal(recoveredPubKey, providedPubKey) {
		return fmt.Errorf("public key mismatch: recovered %x, provided %x", recoveredPubKey, providedPubKey)
	}

	return nil
}

func VerifyEOSCoin(coin, addr, msg, sign, publicKey string) error {
	if publicKey == "" || publicKey == "null" {
		return fmt.Errorf("EOS coin %s missing public key", coin)
	}

	// Convert EOS signature to normal ECDSA format
	normalSig, err := convertEOSSignatureToECDSA(sign)
	if err != nil {
		return fmt.Errorf("failed to convert EOS signature to ECDSA format, coin:%s, addr:%s, error:%v", coin, addr, err)
	}

	// Hash the message using ECDSA format (same as OKX)
	hash := HashEosMsg(OKXMessageSignatureHeader, msg)

	// Recover public key using normal ECDSA
	recoveredPub, err := sigToPub(hash, normalSig)
	if err != nil {
		return fmt.Errorf("failed to recover public key from ECDSA signature, coin:%s, addr:%s, error:%v", coin, addr, err)
	}

	// Convert EOS public key to normal hex format
	expectedPubKeyHex, err := eosPublicKeyToHex(publicKey)
	if err != nil {
		return fmt.Errorf("failed to convert EOS public key to hex, coin:%s, addr:%s, error:%v", coin, addr, err)
	}

	// Convert recovered public key to compressed hex format
	recoveredPubKeyHex := hex.EncodeToString(recoveredPub.SerializeCompressed())

	// Compare the public keys
	if !strings.EqualFold(expectedPubKeyHex, recoveredPubKeyHex) {
		return fmt.Errorf("EOS public key mismatch, coin:%s, expected:%s, recovered:%s", coin, expectedPubKeyHex, recoveredPubKeyHex)
	}

	return nil
}

// convertEOSSignatureToECDSA converts EOS signature format to normal ECDSA format
func convertEOSSignatureToECDSA(eosSignature string) ([]byte, error) {
	if eosSignature == "" {
		return nil, fmt.Errorf("signature cannot be empty")
	}

	if !strings.HasPrefix(eosSignature, "SIG_K1_") {
		return nil, fmt.Errorf("EOS signature must start with 'SIG_K1_'")
	}

	// Remove SIG_K1_ prefix
	base58Part := eosSignature[len("SIG_K1_"):]

	// Decode Base58
	decoded := base58.Decode(base58Part)

	// EOS signature format: 65 bytes signature data + 4 bytes checksum
	if len(decoded) != 69 {
		return nil, fmt.Errorf("signature length incorrect, expected 69 bytes, actual %d bytes", len(decoded))
	}

	// Extract signature part (remove checksum)
	signatureBytes := decoded[:len(decoded)-4]

	// Verify checksum
	// Checksum = RIPEMD160(signatureBytes + "K1")
	ripemd := ripemd160.New()
	ripemd.Write(signatureBytes)
	ripemd.Write([]byte("K1"))
	checksum := ripemd.Sum(nil)
	expectedChecksum := checksum[:4]
	actualChecksum := decoded[len(decoded)-4:]

	if !bytes.Equal(expectedChecksum, actualChecksum) {
		return nil, fmt.Errorf("signature checksum verification failed")
	}

	// Extract recovery ID, r, and s from EOS signature format
	recoveryID := signatureBytes[0]
	r := signatureBytes[1:33]
	s := signatureBytes[33:65]

	// Create normal ECDSA signature format: [r(32) + s(32) + recovery_id(1)]
	normalSig := make([]byte, 65)
	copy(normalSig[0:32], r)
	copy(normalSig[32:64], s)
	normalSig[64] = recoveryID

	return normalSig, nil
}

// eosPublicKeyToHex converts EOS format public key to hex format (compressed)
// Based on the Java implementation provided
func eosPublicKeyToHex(eosPublicKey string) (string, error) {
	const EOSPrefix = "EOS"

	if eosPublicKey == "" {
		return "", fmt.Errorf("public key cannot be empty")
	}

	if !strings.HasPrefix(eosPublicKey, EOSPrefix) {
		return "", fmt.Errorf("EOS public key must start with '%s'", EOSPrefix)
	}

	// Remove EOS prefix
	base58Part := eosPublicKey[len(EOSPrefix):]

	// Decode Base58
	decoded := base58.Decode(base58Part)
	if len(decoded) < 4 {
		return "", fmt.Errorf("invalid decoded length")
	}

	// Extract public key part (remove checksum)
	publicKeyBytes := decoded[:len(decoded)-4]

	// Verify checksum
	ripemd := ripemd160.New()
	ripemd.Write(publicKeyBytes)
	checksum := ripemd.Sum(nil)
	expectedChecksum := checksum[:4]
	actualChecksum := decoded[len(decoded)-4:]

	if !bytes.Equal(expectedChecksum, actualChecksum) {
		return "", fmt.Errorf("public key checksum verification failed")
	}

	// Convert to hex
	return hex.EncodeToString(publicKeyBytes), nil
}

func GenerateICXAddress(publicKey []byte) (string, error) {
	if len(publicKey) < 2 {
		return "", fmt.Errorf("public key too short, length: %d", len(publicKey))
	}

	pub := publicKey[1:]

	hasher := sha3.New256()
	hasher.Write(pub)
	hash := hasher.Sum(nil)

	if len(hash) < 20 {
		return "", fmt.Errorf("hash too short, length: %d", len(hash))
	}
	result := hash[len(hash)-20:]

	addr := "hx" + hex.EncodeToString(result)
	return addr, nil
}

// ConvertETHToLATAddress
func ConvertETHToLATAddress(ethAddress string) (string, error) {
	if len(ethAddress) < 2 || ethAddress[:2] != "0x" {
		return "", fmt.Errorf("invalid ETH address format: %s", ethAddress)
	}

	addressBytes, err := hex.DecodeString(ethAddress[2:])
	if err != nil {
		return "", fmt.Errorf("failed to decode ETH address: %v", err)
	}

	latAddress, err := bech32.EncodeFromBase256("lat", addressBytes)
	if err != nil {
		return "", fmt.Errorf("failed to encode LAT address: %v", err)
	}

	return latAddress, nil
}

func GenerateONEAddress(publicKey []byte) (string, error) {
	if len(publicKey) < 65 {
		return "", fmt.Errorf("public key too short, expected 65 bytes, got %d", len(publicKey))
	}

	var uncompressed []byte
	if publicKey[0] == 0x04 && len(publicKey) == 65 {
		uncompressed = publicKey[1:]
	} else {
		return "", fmt.Errorf("unsupported public key format, expected uncompressed format")
	}

	hasher := sha3.NewLegacyKeccak256()
	hasher.Write(uncompressed)
	hash := hasher.Sum(nil)

	ethAddressBytes := hash[len(hash)-20:]
	oneAddress, err := bech32.EncodeFromBase256("one", ethAddressBytes)
	if err != nil {
		return "", fmt.Errorf("failed to encode ONE address: %v", err)
	}

	return oneAddress, nil
}
