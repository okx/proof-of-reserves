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
	"strings"
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

// GetAdaAddressFromPublicKey generates a Cardano (ADA) address from a public key
// Based on the Java implementation provided - uses EnterpriseKey type (97) for mainnet
func GetAdaAddressFromPublicKey(publicKeyHex string) (string, error) {
	publicKeyBytes, err := Decode(publicKeyHex)
	if err != nil {
		return "", fmt.Errorf("invalid public key: %v", err)
	}

	// Check if public key length is 32 bytes
	if len(publicKeyBytes) != 32 {
		return "", fmt.Errorf("invalid payment key length: expected 32, got %d", len(publicKeyBytes))
	}

	// Create address data: [97, blake2b224(paymentPublicKey)]
	addressData := make([]byte, 0, 29) // 1 byte type + 28 bytes hash

	// Address type: 97 for EnterpriseKey mainnet
	addressData = append(addressData, 97)

	// Blake2b-224 hash of the payment public key
	paymentHash := calculateBlake2b224(publicKeyBytes)
	addressData = append(addressData, paymentHash...)

	// Encode using Bech32 with "addr" human-readable part
	return encodeBech32("addr", addressData)
}

// GetNearAddressFromPublicKey generates a NEAR address from a public key
// Based on the Java implementation: public key hex is the address
func GetNearAddressFromPublicKey(publicKeyHex string) (string, error) {
	// Validate that the public key is valid hex
	_, err := Decode(publicKeyHex)
	if err != nil {
		return "", fmt.Errorf("invalid public key: %v", err)
	}

	// NEAR addresses are just the public key in hex format (without 0x prefix)
	// Remove the 0x prefix if present
	cleanPubKey := strings.TrimPrefix(publicKeyHex, "0x")

	return cleanPubKey, nil
}

// GetHbarAddressFromPublicKey generates a HBAR address from a public key
// Based on the Java implementation: public key string is the address
func GetHbarAddressFromPublicKey(publicKeyHex string) (string, error) {
	// Validate that the public key is valid hex
	_, err := Decode(publicKeyHex)
	if err != nil {
		return "", fmt.Errorf("invalid public key: %v", err)
	}

	// HBAR addresses are just the public key hex string (with 0x prefix)
	// Return the public key as-is
	return publicKeyHex, nil
}

// encodeBech32 encodes data using Bech32 encoding
// Based on the Java implementation provided
func encodeBech32(hrp string, data []byte) (string, error) {
	// Convert bytes to 5-bit groups
	fiveBitData := convertToFiveBitGroups(data)

	// Create checksum
	checksum := createChecksum(hrp, fiveBitData)

	// Combine data and checksum
	combined := make([]byte, len(fiveBitData)+len(checksum))
	copy(combined, fiveBitData)
	copy(combined[len(fiveBitData):], checksum)

	// Build the final string
	var sb strings.Builder
	sb.WriteString(strings.ToLower(hrp))
	sb.WriteByte('1')

	// Map 5-bit values to characters
	charset := "qpzry9x8gf2tvdw0s3jn54khce6mua7l"
	for _, b := range combined {
		if int(b) >= len(charset) {
			return "", fmt.Errorf("invalid 5-bit value: %d", b)
		}
		sb.WriteByte(charset[b])
	}

	return sb.String(), nil
}

// createChecksum creates a Bech32 checksum for the given hrp and data
func createChecksum(hrp string, values []byte) []byte {
	hrpExpanded := expandHrp(hrp)
	enc := make([]byte, len(hrpExpanded)+len(values)+6)
	copy(enc, hrpExpanded)
	copy(enc[len(hrpExpanded):], values)

	mod := polymod(enc) ^ 1
	ret := make([]byte, 6)

	for i := 0; i < 6; i++ {
		ret[i] = byte((mod >> (5 * (5 - i))) & 31)
	}

	return ret
}

// expandHrp expands the human-readable part for checksum calculation
func expandHrp(hrp string) []byte {
	hrpLength := len(hrp)
	ret := make([]byte, hrpLength*2+1)

	for i := 0; i < hrpLength; i++ {
		c := int(hrp[i]) & 127
		ret[i] = byte((c >> 5) & 7)
		ret[i+hrpLength+1] = byte(c & 31)
	}

	ret[hrpLength] = 0
	return ret
}

// polymod performs the polynomial modulo operation for Bech32 checksum
func polymod(values []byte) int {
	c := 1
	for _, v := range values {
		c0 := (c >> 25) & 255
		c = ((c & 33554431) << 5) ^ int(v&255)

		if (c0 & 1) != 0 {
			c ^= 996825010
		}
		if (c0 & 2) != 0 {
			c ^= 642813549
		}
		if (c0 & 4) != 0 {
			c ^= 513874426
		}
		if (c0 & 8) != 0 {
			c ^= 1027748829
		}
		if (c0 & 16) != 0 {
			c ^= 705979059
		}
	}

	return c
}

// convertToFiveBitGroups converts byte array to 5-bit groups
// Based on the Java encodeCombineData implementation
func convertToFiveBitGroups(values []byte) []byte {
	var split []byte

	// Calculate size: if values.length * 8 % 5 == 0, then size = values.length * 8 / 5
	// else size = values.length * 8 / 5 + 1
	if len(values)*8%5 == 0 {
		split = make([]byte, len(values)*8/5)
	} else {
		split = make([]byte, len(values)*8/5+1)
	}

	// Convert bytes to 5-bit groups
	for i := 0; i < len(split)-1; i++ {
		switch i % 8 {
		case 0:
			split[i] = byte((values[i/8*5+0] & 248) >> 3)
		case 1:
			if i/8*5+1 < len(values) {
				split[i] = byte((values[i/8*5+0]&7)<<2 | (values[i/8*5+1]&192)>>6)
			} else {
				split[i] = byte((values[i/8*5+0] & 7) << 2)
			}
		case 2:
			if i/8*5+1 < len(values) {
				split[i] = byte((values[i/8*5+1] & 62) >> 1)
			}
		case 3:
			if i/8*5+2 < len(values) {
				split[i] = byte((values[i/8*5+1]&1)<<4 | (values[i/8*5+2]&240)>>4)
			} else {
				split[i] = byte((values[i/8*5+1] & 1) << 4)
			}
		case 4:
			if i/8*5+3 < len(values) {
				split[i] = byte((values[i/8*5+2]&15)<<1 | (values[i/8*5+3]&128)>>7)
			} else {
				split[i] = byte((values[i/8*5+2] & 15) << 1)
			}
		case 5:
			if i/8*5+3 < len(values) {
				split[i] = byte((values[i/8*5+3] & 124) >> 2)
			}
		case 6:
			if i/8*5+4 < len(values) {
				split[i] = byte((values[i/8*5+3]&3)<<3 | (values[i/8*5+4]&224)>>5)
			} else {
				split[i] = byte((values[i/8*5+3] & 3) << 3)
			}
		case 7:
			if i/8*5+4 < len(values) {
				split[i] = byte(values[i/8*5+4] & 31)
			}
		}
	}

	// Handle the last element
	i := len(split) - 1
	if i >= 0 {
		switch i % 8 {
		case 1:
			if i/8*5+0 < len(values) {
				split[i] = byte((values[i/8*5+0] & 7) << 2)
			}
		case 3:
			if i/8*5+1 < len(values) {
				split[i] = byte((values[i/8*5+1] & 1) << 4)
			}
		case 4:
			if i/8*5+2 < len(values) {
				split[i] = byte((values[i/8*5+2] & 15) << 1)
			}
		case 6:
			if i/8*5+3 < len(values) {
				split[i] = byte((values[i/8*5+3] & 3) << 3)
			}
		case 7:
			if i/8*5+4 < len(values) {
				split[i] = byte(values[i/8*5+4] & 31)
			}
		}
	}

	return split
}

// GetXlmAddressFromPublicKey generates a Stellar (XLM) address from a public key
// Based on OKX go-wallet-sdk implementation
func GetXlmAddressFromPublicKey(publicKeyHex string) (string, error) {
	publicKeyBytes, err := Decode(publicKeyHex)
	if err != nil {
		return "", fmt.Errorf("invalid public key: %v", err)
	}

	// XLM addresses use ed25519 public keys directly
	// Version byte for account ID is 6 << 3 = 48 (0x30)
	versionByte := byte(48)

	// Combine version byte and public key
	versionedKey := append([]byte{versionByte}, publicKeyBytes...)

	// Calculate CRC16 checksum (little endian)
	checksum := calculateCRC16(versionedKey)

	// Combine versioned key and checksum
	addressBytes := append(versionedKey, checksum...)

	// Base32 encode with standard alphabet (ABCDEFGHIJKLMNOPQRSTUVWXYZ234567)
	// and no padding
	encoding := base32.StdEncoding.WithPadding(base32.NoPadding)

	return encoding.EncodeToString(addressBytes), nil
}

// calculateCRC16 calculates CRC16 checksum for Stellar address
// This is a simplified implementation - in production you might want to use a proper CRC16 library
func calculateCRC16(data []byte) []byte {
	// This is a basic CRC16 implementation
	// For production use, consider using a proper CRC16 library like github.com/okx/go-wallet-sdk/coins/stellar/strkey/internal/crc16
	crc := uint16(0)
	for _, b := range data {
		crc ^= uint16(b) << 8
		for i := 0; i < 8; i++ {
			if crc&0x8000 != 0 {
				crc = (crc << 1) ^ 0x1021
			} else {
				crc <<= 1
			}
		}
	}

	// Convert to little endian bytes
	result := make([]byte, 2)
	result[0] = byte(crc & 0xFF)
	result[1] = byte((crc >> 8) & 0xFF)

	return result
}

// calculateBlake2b224 calculates Blake2b-224 hash
// Based on the Java implementation that accepts multiple byte arrays
func calculateBlake2b224(data ...[]byte) []byte {
	// Use the official Go blake2b implementation with 28-byte output (224 bits)
	hasher, err := blake2b.New(&blake2b.Config{Size: 28})
	if err != nil {
		// Fallback to SHA256 if Blake2b is not available
		// Combine all input arrays and hash
		var combined []byte
		for _, d := range data {
			combined = append(combined, d...)
		}
		hash := sha256.Sum256(combined)
		return hash[:28]
	}

	// Write all input byte arrays to the hasher
	for _, d := range data {
		hasher.Write(d)
	}

	return hasher.Sum(nil)
}
