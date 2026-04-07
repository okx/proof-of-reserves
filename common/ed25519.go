package common

import (
	"errors"
	"github.com/dchest/blake2b"
	"github.com/martinboehm/btcutil/base58"
)

func GetDotAddressFromPublicKey(publicKeyHex string) (string, error) {
	return GetSubstrateAddressFromPublicKey(publicKeyHex, 0)
}

func GetSubstrateAddressFromPublicKey(publicKeyHex string, network uint16) (string, error) {
	publicKeyBytes, _ := Decode(publicKeyHex)
	if len(publicKeyBytes) != 32 {
		return "", errors.New("public hash length is not equal 32")
	}

	ssPrefix := []byte("SS58PRE")

	var prefix []byte
	if network < 64 {
		prefix = []byte{byte(network)}
	} else {
		first := byte(((network & 0xfc) >> 2) | 0x40)
		second := byte((network >> 8) | ((network & 0x03) << 6))
		prefix = []byte{first, second}
	}

	payload := appendBytes(prefix, publicKeyBytes)
	input := appendBytes(ssPrefix, payload)
	ck := blake2b.Sum512(input)
	checksum := ck[:2]
	address := base58.Encode(appendBytes(payload, checksum))
	if address == "" {
		return address, errors.New("base58 encode error")
	}
	return address, nil
}

func appendBytes(data1, data2 []byte) []byte {
	if data2 == nil {
		return data1
	}
	return append(data1, data2...)
}
