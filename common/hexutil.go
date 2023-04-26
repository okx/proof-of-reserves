package common

import (
	"encoding/hex"
	"golang.org/x/crypto/sha3"
	"strconv"
)

// Errors
var (
	ErrEmptyString = &decError{"empty hex string"}
	ErrSyntax      = &decError{"invalid hex string"}
	ErrOddLength   = &decError{"hex string of odd length"}
	ErrUint64Range = &decError{"hex number > 64 bits"}
)

type decError struct{ msg string }

func (err decError) Error() string { return err.msg }

// Decode decodes a hex string with 0x prefix.
func Decode(input string) ([]byte, error) {
	if len(input) == 0 {
		return nil, ErrEmptyString
	}
	var b []byte
	var err error
	if has0xPrefix(input) {
		b, err = hex.DecodeString(input[2:])
	} else {
		b, err = hex.DecodeString(input)
	}
	if err != nil {
		err = mapError(err)
	}
	return b, err

}

// MustDecode decodes a hex string with 0x prefix. It panics for invalid input.
func MustDecode(input string) []byte {
	dec, err := Decode(input)
	if err != nil {
		panic(err)
	}
	return dec
}

// Encode encodes b as a hex string with 0x prefix.
func Encode(b []byte) string {
	enc := make([]byte, len(b)*2+2)
	copy(enc, "0x")
	hex.Encode(enc[2:], b)
	return string(enc)
}

func mapError(err error) error {
	if err, ok := err.(*strconv.NumError); ok {
		switch err.Err {
		case strconv.ErrRange:
			return ErrUint64Range
		case strconv.ErrSyntax:
			return ErrSyntax
		}
	}
	if _, ok := err.(hex.InvalidByteError); ok {
		return ErrSyntax
	}
	if err == hex.ErrLength {
		return ErrOddLength
	}
	return err
}

func Sha256Hash(bytes []byte) []byte {
	sha256 := sha3.New256()
	sha256.Write(bytes)
	return sha256.Sum(nil)
}
