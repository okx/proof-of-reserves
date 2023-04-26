package common

import (
	"errors"
	"github.com/dchest/blake2b"
	"github.com/martinboehm/btcutil/base58"
)

func GetDotAddressFromPublicKey(publicKeyHex string) (string, error) {
	prefix := []byte{0x00}
	ssPrefix := []byte{0x53, 0x53, 0x35, 0x38, 0x50, 0x52, 0x45}
	publicKeyBytes, _ := Decode(publicKeyHex)
	if len(publicKeyBytes) != 32 {
		return "", errors.New("public hash length is not equal 32")
	}
	payload := appendBytes(prefix, publicKeyBytes)
	input := appendBytes(ssPrefix, payload)
	ck := blake2b.Sum512(input)
	checkum := ck[:2]
	address := base58.Encode(appendBytes(payload, checkum))
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
