package common

import (
	"bytes"
	"crypto/sha256"
	"fmt"

	"github.com/btcsuite/btcd/chaincfg/chainhash"
	"github.com/btcsuite/btcd/wire"
)

func GetNetWork() []byte {
	return []byte{0x41}
}

func HashUtxoCoinTypeMsg(msgHeader, msg string) []byte {
	var buf bytes.Buffer
	wire.WriteVarString(&buf, 0, msgHeader)
	wire.WriteVarString(&buf, 0, msg)
	expectedMessageHash := chainhash.DoubleHashB(buf.Bytes())
	return expectedMessageHash
}

func HashEd25519Msg(msgHeader, msg string) []byte {
	var buf bytes.Buffer
	wire.WriteVarString(&buf, 0, msgHeader)
	wire.WriteVarString(&buf, 0, msg)
	return buf.Bytes()
}

func HashEcdsaMsg(msgHeader, msg string) []byte {
	var buf bytes.Buffer
	wire.WriteVarString(&buf, 0, msgHeader)
	wire.WriteVarString(&buf, 0, msg)
	return Keccak256(buf.Bytes())
}

func HashEosMsg(msgHeader, msg string) []byte {
	var buf bytes.Buffer
	wire.WriteVarString(&buf, 0, msgHeader)
	wire.WriteVarString(&buf, 0, msg)

	hash := sha256.Sum256(buf.Bytes())
	return hash[:]
}

func HashEvmCoinTypeMsg(msgHeader, msg string) []byte {
	var buf bytes.Buffer
	b := Keccak256([]byte(msg))
	buf.Write([]byte(msgHeader))
	buf.Write(b)
	expectedMessageHash := Keccak256(buf.Bytes())
	return expectedMessageHash
}

func HashTrxMsg(msg string) []byte {
	var buf bytes.Buffer
	b := Keccak256([]byte(msg))
	buf.Write([]byte(TronMessageSignatureHeader))
	buf.Write(b)
	expectedMessageHash := Keccak256(buf.Bytes())
	return expectedMessageHash
}

func HashTrxMsgV2(msg string) []byte {
	length := fmt.Sprintf("%d", len(msg))

	var buf bytes.Buffer
	buf.WriteString(TronMessageV2SignatureHeader)
	buf.WriteString(length)
	buf.WriteString(msg)

	expectedMessageHash := Keccak256(buf.Bytes())
	return expectedMessageHash
}
