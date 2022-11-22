package common

import (
	"bytes"
	"github.com/btcsuite/btcd/chaincfg/chainhash"
	"github.com/btcsuite/btcd/wire"
)

const btcMessageSignatureHeader = "Bitcoin Signed Message:\n"
const solMessageSignatureHeader = "OKX Signed Message:\n"

var (
	ethMessageSignatureHeader = []byte("\x19Ethereum Signed Message:\n32")
	trxMessageSignatureHeader = []byte("\x19TRON Signed Message:\n32")
)

func GetNetWork() []byte {
	return []byte{0x41}
}

func HashBTCMsg(msg string) []byte {
	var buf bytes.Buffer
	wire.WriteVarString(&buf, 0, btcMessageSignatureHeader)
	wire.WriteVarString(&buf, 0, msg)
	expectedMessageHash := chainhash.DoubleHashB(buf.Bytes())
	return expectedMessageHash
}

func SolMsg(msg string) []byte {
	var buf bytes.Buffer
	wire.WriteVarString(&buf, 0, solMessageSignatureHeader)
	wire.WriteVarString(&buf, 0, msg)
	return buf.Bytes()
}

func HashETHMsg(msg string) []byte {
	var buf bytes.Buffer
	b := Keccak256([]byte(msg))
	buf.Write(ethMessageSignatureHeader)
	buf.Write(b)
	expectedMessageHash := Keccak256(buf.Bytes())
	return expectedMessageHash
}

func HashTrxMsg(msg string) []byte {
	var buf bytes.Buffer
	b := Keccak256([]byte(msg))
	buf.Write(trxMessageSignatureHeader)
	buf.Write(b)
	expectedMessageHash := Keccak256(buf.Bytes())
	return expectedMessageHash
}
