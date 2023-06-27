package common

import (
	"bytes"
	"crypto/sha256"
	"encoding/base32"
	"encoding/binary"
	"encoding/hex"
	"fmt"
	"github.com/dchest/blake2b"
	"github.com/filecoin-project/go-address"
	builtintypes "github.com/filecoin-project/go-state-types/builtin"
	"github.com/martinboehm/btcutil/base58"
)

var (
	payloadHashConfig  = &blake2b.Config{Size: 20}
	checksumHashConfig = &blake2b.Config{Size: 4}

	AddressEncoding = base32.NewEncoding(encodeStd)
)

const (
	encodeStd = "abcdefghijklmnopqrstuvwxyz234567"
)

func GetFilAddressFromPublicKey(publicKeyHex string) string {
	publicKeyBytes, _ := Decode(publicKeyHex)
	pubKeyHash := hash_cal(publicKeyBytes, payloadHashConfig)

	explen := 1 + len(pubKeyHash)
	buf := make([]byte, explen)
	var protocol byte = 1
	buf[0] = protocol
	copy(buf[1:], pubKeyHash)

	cksm := hash_cal(buf, checksumHashConfig)
	address := "f" + fmt.Sprintf("%d", protocol) + AddressEncoding.WithPadding(-1).EncodeToString(append(pubKeyHash, cksm[:]...))

	return address
}

func hash_cal(ingest []byte, cfg *blake2b.Config) []byte {
	hasher, err := blake2b.New(cfg)
	if err != nil {
		// If this happens sth is very wrong.
		panic(fmt.Sprintf("invalid address hash configuration: %v", err)) // ok
	}
	if _, err := hasher.Write(ingest); err != nil {
		// blake2bs Write implementation never returns an error in its current
		// setup. So if this happens sth went very wrong.
		panic(fmt.Sprintf("blake2b is unable to process hashes: %v", err)) // ok
	}
	return hasher.Sum(nil)
}

func GetElfAddressFromPublicKey(publicKeyBytes []byte) string {
	firstBytes := sha256.Sum256(publicKeyBytes)
	secondBytes := sha256.Sum256(firstBytes[:])

	return encodeCheck(secondBytes[:])
}

func encodeCheck(input []byte) string {
	b := make([]byte, 0, 1+len(input)+4)
	b = append(b, input[:]...)
	cksum := checksum(b)
	b = append(b, cksum[:]...)

	return base58.Encode(b)
}

func checksum(input []byte) (cksum [4]byte) {
	h := sha256.Sum256(input)
	h2 := sha256.Sum256(h[:])
	copy(cksum[:], h2[:4])

	return
}

var maskedIDPrefix = [20 - 8]byte{0xff}

func ConvertEthAddressToFilecoinAddress(ethHash []byte) (address.Address, error) {
	address.CurrentNetwork = address.Mainnet
	if bytes.HasPrefix(ethHash[:], maskedIDPrefix[:]) {
		// This is a masked ID address.
		id := binary.BigEndian.Uint64(ethHash[12:])
		return address.NewIDAddress(id)
	}

	// Otherwise, translate the address into an address controlled by the
	// Ethereum Address Manager.
	addr, err := address.NewDelegatedAddress(builtintypes.EthereumAddressManagerActorID, ethHash[:])
	if err != nil {
		return address.Undef, fmt.Errorf("failed to translate supplied address (%s) into a "+
			"Filecoin f4 address: %w", hex.EncodeToString(ethHash[:]), err)
	}
	return addr, nil
}
