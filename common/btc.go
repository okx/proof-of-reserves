package common

import (
	"encoding/base64"
	"encoding/hex"
	"fmt"
	secp_ecdsa "github.com/decred/dcrd/dcrec/secp256k1/v4/ecdsa"
	"regexp"
	"strings"
)

// CreateAddressDescriptor generate address descriptor
// https://github.com/bitcoin/bitcoin/blob/master/doc/descriptors.md
func CreateAddressDescriptor(addrType, redeemScript string, mSigns, nKeys int) (descriptor string, err error) {
	var orderedPubKeys []string
	var pubKeyLen int
	pubKeyLen = 33 * 2
	if addrType == "P2PKH" {
		orderedPubKeys = append(orderedPubKeys, redeemScript)
	} else {
		redeemScript = redeemScript[2:]
		for i := 0; i < nKeys; i++ {
			pubKey := redeemScript[2 : 2+pubKeyLen]
			orderedPubKeys = append(orderedPubKeys, pubKey)
			redeemScript = redeemScript[2+pubKeyLen:]
		}
	}
	orderedPubKeysJoin := strings.Join(orderedPubKeys, ",")
	switch addrType {
	case "P2PKH":
		descriptor = fmt.Sprintf("pkh(%s)", orderedPubKeysJoin)
	case "P2SH":
		descriptor = fmt.Sprintf("sh(multi(%d,%s))", mSigns, orderedPubKeysJoin)
	case "P2WSH":
		descriptor = fmt.Sprintf("wsh(multi(%d,%s))", mSigns, orderedPubKeysJoin)
	default:
		descriptor = fmt.Sprintf("sh(multi(%d,%s))", mSigns, orderedPubKeysJoin)
	}

	return descriptor, nil
}

func GuessAddressType(address string) string {
	match1, _ := regexp.MatchString("^[1-9A-Za-z]{26,35}$", address)
	if match1 {
		if address[0:1] == "1" {
			return "P2PKH"
		}
		if address[0:1] == "3" {
			return "P2SH"
		}
	}
	if len(address) == 40 {
		return "P2WPKH"
	}

	if len(address) == 64 {
		return "P2WSH"
	}
	match2, _ := regexp.MatchString("^bc1[0-9a-zA-Z]{11,71}$", strings.ToLower(address))
	if match2 {
		return "P2WSH"
	}

	return ""
}

func RecoveryPubKeyFromSign(address, msg, sign string) (pubKey string) {
	hash := HashBTCMsg(msg)
	b, err := base64.StdEncoding.DecodeString(sign)
	if err != nil {
		return ""
	}
	pub, ok, err := secp_ecdsa.RecoverCompact(b, hash)
	if err != nil || !ok || pub == nil {
		return ""
	}
	return hex.EncodeToString(pub.SerializeCompressed())
}
