package common

import (
	"crypto/ed25519"
	"crypto/sha256"
	"encoding/base64"
	"errors"
	"fmt"
	"github.com/btcsuite/btcd/btcutil"
	"github.com/btcsuite/btcd/btcutil/base58"
	"github.com/btcsuite/btcd/chaincfg"
	"github.com/btcsuite/btcd/txscript"
	secp_ecdsa "github.com/decred/dcrd/dcrec/secp256k1/v4/ecdsa"
	"golang.org/x/crypto/sha3"
)

var (
	ErrInvalidAddr = errors.New("invalid address")
	ErrInvalidSign = errors.New("can't verify signature")
)

func VerifySol(pubkey, msg, sign string) error {
	hash := SolMsg(msg)
	res := base58.Decode(sign)
	pub := base58.Decode(pubkey)
	if ok := ed25519.Verify(pub, hash, res); !ok {
		return ErrInvalidSign
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
	newAddr := base58.CheckEncode(newHash, GetNetWork()[0])
	if addr != newAddr {
		return ErrInvalidSign
	}
	return nil
}

func VerifyBTCWitness(addr, msg, pkScript, sign1, sign2 string) bool {
	script := MustDecode(pkScript)
	h := sha256.New()
	h.Write(script)
	witnessProg := h.Sum(nil)
	addressWitnessScriptHash, err := btcutil.NewAddressWitnessScriptHash(witnessProg, &chaincfg.MainNetParams)
	if err != nil {
		return false
	}
	if addressWitnessScriptHash.EncodeAddress() != addr {
		return false
	}
	addr1, err := SigToAddrBTC(msg, sign1)
	if err != nil {
		return false
	}
	addr2, err := SigToAddrBTC(msg, sign2)
	if err != nil {
		return false
	}
	m := map[string]struct{}{addr1: {}, addr2: {}}
	_, addrs, _, err := txscript.ExtractPkScriptAddrs(script, &chaincfg.MainNetParams)
	if err != nil {
		return false
	}
	for _, v := range addrs {
		delete(m, v.EncodeAddress())
	}
	if len(m) > 1 {
		return false
	}
	return true
}

func VerifyBTCP2PKH(addr, msg, pkScript, sign1, sign2 string) bool {
	fmt.Println("addr", addr)
	script := MustDecode(pkScript)
	addrPub, err := btcutil.NewAddressPubKeyHash(script[3:23], &chaincfg.MainNetParams)
	if err != nil {
		return false
	}
	if addrPub.EncodeAddress() != addr {
		return false
	}
	addr1, err := SigToAddrBTC(msg, sign1)
	if err != nil {
		return false
	}
	addr2, err := SigToAddrBTC(msg, sign2)
	if err != nil {
		return false
	}
	return addr2 == addr || addr1 == addr
}

func VerifyBTC(addr, msg, sign string) error {
	hash := HashBTCMsg(msg)
	b, err := base64.StdEncoding.DecodeString(sign)
	if err != nil {
		return ErrInvalidSign
	}
	//{"address":"1GhLyRg4zzFixW3ZY5ViFzT4W5zTT9h7Pc","message":"hello world","signature1":"Hz+cZI5GfSzNSvBpna20diV47/rhlQMRQTNGZd9sI4UZQaWH4ZY3KJA4IlcP5bwuicO+myA4vLdiMkj7OU+rDpg="}
	pub, ok, err := secp_ecdsa.RecoverCompact(b, hash)
	if err != nil || !ok || pub == nil {
		return ErrInvalidSign
	}
	if _, err := btcutil.DecodeAddress(addr, &chaincfg.MainNetParams); err != nil {
		return ErrInvalidSign
	}
	addrPub, err := btcutil.NewAddressPubKey(pub.SerializeCompressed(), &chaincfg.MainNetParams)
	if err != nil || addrPub.EncodeAddress() != addr {
		return ErrInvalidSign
	}
	return nil
}

func SigToAddrBTC(msg, sign string) (string, error) {
	hash := HashBTCMsg(msg)
	b, err := base64.StdEncoding.DecodeString(sign)
	if err != nil {
		return "", ErrInvalidSign
	}
	pub, ok, err := secp_ecdsa.RecoverCompact(b, hash)
	if err != nil || !ok || pub == nil {
		return "", ErrInvalidSign
	}
	addrPub, err := btcutil.NewAddressPubKey(pub.SerializeCompressed(), &chaincfg.MainNetParams)
	if err != nil {
		return "", ErrInvalidSign
	}
	return addrPub.EncodeAddress(), nil
}

func VerifyETH(addr, msg, sign string) error {
	hash := HashETHMsg(msg)
	s := MustDecode(sign)
	pub, err := SigToPub(hash, s)
	if err != nil {
		return ErrInvalidAddr
	}
	if HexToAddress(addr) != PubkeyToAddress(*pub) {
		fmt.Println(PubkeyToAddress(*pub).String())
		return ErrInvalidSign
	}
	if !VerifySignAddr(HexToAddress(addr), hash, s) {
		return ErrInvalidSign
	}
	return nil
}
