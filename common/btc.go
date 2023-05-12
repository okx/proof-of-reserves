package common

import (
	"encoding/base64"
	"encoding/hex"
	"fmt"
	secp_ecdsa "github.com/decred/dcrd/dcrec/secp256k1/v4/ecdsa"
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

func RecoveryPubKeyFromSign(address, msg, sign string) (pubKey string) {
	hash := HashUtxoCoinTypeMsg(BtcMessageSignatureHeader, msg)
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
