package common

import (
	"crypto/ed25519"
	"crypto/sha256"
	"encoding/base64"
	"encoding/hex"
	"errors"
	"fmt"
	"github.com/Conflux-Chain/go-conflux-sdk/types/cfxaddress"
	"github.com/martinboehm/btcd/txscript"
	"github.com/martinboehm/btcutil"
	"github.com/martinboehm/btcutil/base58"
	"github.com/martinboehm/btcutil/bech32"
	"github.com/martinboehm/btcutil/chaincfg"
	"golang.org/x/crypto/ripemd160"
	"strings"

	secp_ecdsa "github.com/decred/dcrd/dcrec/secp256k1/v4/ecdsa"
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
		return errors.New(fmt.Sprintf("unexpected failure, failed to deserialize pubkey (%x): %v", p[:], err))
	}
	var sig Signature
	if err := sig.Deserialize(&s); err != nil {
		return errors.New(fmt.Sprintf("unexpected failure, failed to deserialize signature (%x): %v", s[:], err))
	}
	res := Verify(&pub, h[:], &sig)
	if !res {
		return errors.New("unexpected failure, failed to verify signature")
	}
	return nil
}

func VerifyTRX(addr, msg, sign string) error {
	hash := HashTrxMsg(msg)
	s := MustDecode(sign)
	pub, err := sigToPub(hash, s)
	if err != nil {
		return ErrInvalidSign
	}
	pubKey := pub.SerializeUncompressed()
	h := sha3.NewLegacyKeccak256()
	h.Write(pubKey[1:])
	newHash := h.Sum(nil)[12:]
	newAddr := base58.CheckEncode(newHash, GetNetWork(), base58.Sha256D)
	if addr != newAddr {
		return ErrInvalidSign
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
		return nil, ErrInvalidSign
	}
	pub, ok, err := secp_ecdsa.RecoverCompact(b, hash)
	if err != nil || !ok || pub == nil {
		return nil, ErrInvalidSign
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
		return ErrInvalidSign
	}
	addrType := GuessUtxoCoinAddressType(addr)
	switch addrType {
	case "P2PKH":
		addrPub, err := btcutil.NewAddressPubKey(pub1, mainNetParams)
		if err != nil || addrPub.EncodeAddress() != addr {
			return fmt.Errorf("address not match,coin: %s, addr: %s, recoverAddr: %s", coin, addr, addrPub.EncodeAddress())
		}
	case "P2SH":
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
		if len(pubs) != 3 {
			return fmt.Errorf("script address pubs num not match, coin:%s, addr:%s, srcNum: %d, num: %d", coin, addr, 3, len(pubs))
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
		if len(pubs) != 3 {
			return fmt.Errorf("script address pubs num not match, coin:%s, addr:%s, srcNum: %d, num: %d", coin, addr, 3, len(pubs))
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
	case "ETH":
		if !VerifySignAddr(HexToAddress(addr), hash, s) {
			return ErrInvalidSign
		}
	}

	if strings.ToLower(addr) != strings.ToLower(recoverAddr) {
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
		return ErrInvalidSign
	}
	var recoverAddr string
	switch coin {
	case "SOL":
		out := [32]byte{}
		byteCount := len(pubkeyBytes)
		if byteCount == 0 {
			return ErrInvalidSign
		}
		max := 32
		if byteCount < max {
			max = byteCount
		}
		copy(out[:], pubkeyBytes[0:max])
		recoverAddr = base58.Encode(out[:])
	case "APTOS":
		publicKey := append(pubkeyBytes, 0x0)
		recoverAddr = "0x" + hex.EncodeToString(Sha256Hash(publicKey))
	case "TONCOIN-NEW":
		a, err := tonWallet.AddressFromPubKey(pubkeyBytes, tonWallet.V3, tonWallet.DefaultSubwallet)
		if err != nil {
			return fmt.Errorf("%s, coin: %s, addr: %s, error: %v", ErrInvalidSign, coin, addr, err)
		}
		recoverAddr = a.String()
	case "DOT":
		rAddr, err := GetDotAddressFromPublicKey(pubkey)
		if err != nil {
			return fmt.Errorf("%s, coin: %s, addr: %s, error: %v", ErrInvalidSign, coin, addr, err)
		}
		recoverAddr = rAddr
	}
	if strings.ToLower(recoverAddr) != strings.ToLower(addr) {
		return fmt.Errorf("recovery address not match, coin:%s, recoverAddr:%s, addr:%s", coin, recoverAddr, addr)
	}
	return nil
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
		return ErrInvalidSign
	}
	pubKey := pub.SerializeUncompressed()

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
	}
	if strings.ToLower(recoverAddr) != strings.ToLower(addr) {
		return fmt.Errorf("recovery address not match, coin:%s, recoverAddr:%s, addr:%s", coin, recoverAddr, addr)
	}

	return nil
}
