package common

// Pearl (PRL) PoR verification for Taproot script-path 2-of-2 addresses.
//
// Pearl's verifymessage is a BIP-322 "simple" variant (pearl-research-labs/pearl
// v1.5.3 wallet/rpc/legacyrpc/bip322.go, verified 2026-10-29): both virtual
// transactions are version 2, and to_sign commits the raw message in its
// OP_RETURN output instead of BIP-322's tagged message hash. OKX addresses are
// 2-of-2 tapscripts, so the report carries both Schnorr signatures plus the
// tapscript, and the proof is checked by running the script engine on the
// witness [signature2, signature1, tapscript, controlBlock], as verifymessage
// does.

import (
	"bytes"
	"encoding/hex"
	"errors"
	"fmt"
	"strings"

	btcec_schnorr "github.com/btcsuite/btcd/btcec/v2/schnorr"
	"github.com/btcsuite/btcd/btcutil/bech32"
	"github.com/btcsuite/btcd/chaincfg/chainhash"
	btc_txscript "github.com/btcsuite/btcd/txscript"
	"github.com/btcsuite/btcd/wire"
)

var (
	ErrPearlScript    = errors.New("invalid Pearl 2-of-2 tapscript")
	ErrPearlSignature = errors.New("can't verify Pearl signature")
)

const (
	// Pearl mainnet bech32m HRP, node/chaincfg/params.go v1.5.3, verified 2026-10-29.
	pearlMainnetHRP   = "prl"
	pearlTapscriptLen = 70
)

// BIP341 NUMS internal key, identical on every network.
// https://github.com/bitcoin/bips/blob/master/bip-0341.mediawiki
var pearlNUMSKey = MustDecode("50929b74c1a04954b78b4b6035e97a5e078a5a0f28ec96d547bfee9ace803ac0")

// VerifyPearl verifies a Pearl PoR proof for a mainnet 2-of-2 address: script is
// the tapscript "20<key1>ac20<key2>ba529c", sign1/sign2 are the hex BIP-340
// signatures of key1/key2.
//
// The address is rebuilt from the tapscript, so a tapscript for any other keys
// cannot pass, and the signatures are checked by the script engine rather than
// against a precomputed hash.
func VerifyPearl(addr, msg, sign1, sign2, script string) error {
	tapscript, err := pearlTapscript(script)
	if err != nil {
		return err
	}
	leaf := btc_txscript.NewBaseTapLeaf(tapscript)
	root := leaf.TapHash()
	internalKey, err := btcec_schnorr.ParsePubKey(pearlNUMSKey)
	if err != nil {
		return fmt.Errorf("failed to parse NUMS key: %v", err)
	}
	program := btcec_schnorr.SerializePubKey(btc_txscript.ComputeTaprootOutputKey(internalKey, root[:]))
	expected, err := pearlAddress(program)
	if err != nil {
		return err
	}
	if !strings.EqualFold(expected, addr) {
		return fmt.Errorf("address not match, coin:PRL, addr:%s, scriptAddr:%s", addr, expected)
	}

	proof := btc_txscript.TapscriptProof{TapLeaf: leaf, RootNode: leaf}
	control := proof.ToControlBlock(internalKey)
	controlBlock, err := control.ToBytes()
	if err != nil {
		return fmt.Errorf("failed to encode control block, addr:%s, error:%v", addr, err)
	}
	sig1, err := pearlSignature(sign1, addr)
	if err != nil {
		return err
	}
	sig2, err := pearlSignature(sign2, addr)
	if err != nil {
		return err
	}

	pkScript := append([]byte{btc_txscript.OP_1, btc_txscript.OP_DATA_32}, program...)
	toSign, err := pearlToSign(msg, pkScript)
	if err != nil {
		return fmt.Errorf("%w: addr:%s, error:%v", ErrPearlSignature, addr, err)
	}
	toSign.TxIn[0].Witness = wire.TxWitness{sig2, sig1, tapscript, controlBlock}

	prevOutFetcher := btc_txscript.NewCannedPrevOutputFetcher(pkScript, 0)
	vm, err := btc_txscript.NewEngine(pkScript, toSign, 0, btc_txscript.StandardVerifyFlags, nil,
		btc_txscript.NewTxSigHashes(toSign, prevOutFetcher), 0, prevOutFetcher)
	if err != nil {
		return fmt.Errorf("%w: addr:%s, error:%v", ErrPearlSignature, addr, err)
	}
	if err := vm.Execute(); err != nil {
		return fmt.Errorf("%w: addr:%s, error:%v", ErrPearlSignature, addr, err)
	}
	return nil
}

// pearlTapscript decodes script and accepts only the canonical 2-of-2 leaf
// PUSH32 key1 CHECKSIG PUSH32 key2 CHECKSIGADD 2 NUMEQUAL with two distinct keys.
func pearlTapscript(script string) ([]byte, error) {
	b, err := hex.DecodeString(strings.TrimPrefix(strings.TrimPrefix(script, "0x"), "0X"))
	if err != nil || len(b) != pearlTapscriptLen ||
		b[0] != btc_txscript.OP_DATA_32 || b[33] != btc_txscript.OP_CHECKSIG ||
		b[34] != btc_txscript.OP_DATA_32 || b[67] != btc_txscript.OP_CHECKSIGADD ||
		b[68] != btc_txscript.OP_2 || b[69] != btc_txscript.OP_NUMEQUAL {
		return nil, fmt.Errorf("%w: %s", ErrPearlScript, script)
	}
	key1, key2 := b[1:33], b[35:67]
	if bytes.Equal(key1, key2) {
		return nil, fmt.Errorf("%w: duplicate keys", ErrPearlScript)
	}
	for _, key := range [][]byte{key1, key2} {
		if _, err := btcec_schnorr.ParsePubKey(key); err != nil {
			return nil, fmt.Errorf("%w: invalid key %x", ErrPearlScript, key)
		}
	}
	return b, nil
}

// pearlAddress encodes a witness v1 program as a Pearl mainnet bech32m address.
func pearlAddress(program []byte) (string, error) {
	converted, err := bech32.ConvertBits(program, 8, 5, true)
	if err != nil {
		return "", fmt.Errorf("failed to convert witness program: %v", err)
	}
	return bech32.EncodeM(pearlMainnetHRP, append([]byte{1}, converted...))
}

// pearlSignature decodes one hex BIP-340 signature: 64 bytes for
// SIGHASH_DEFAULT, or 65 bytes carrying an explicit sighash type.
func pearlSignature(sign, addr string) ([]byte, error) {
	if sign == "" || sign == "null" || sign == "\\N" {
		return nil, fmt.Errorf("%w: both signatures are required, addr:%s", ErrPearlSignature, addr)
	}
	b, err := hex.DecodeString(strings.TrimPrefix(strings.TrimPrefix(sign, "0x"), "0X"))
	if err != nil || (len(b) != 64 && len(b) != 65) {
		return nil, fmt.Errorf("%w: invalid signature encoding, addr:%s", ErrPearlSignature, addr)
	}
	return b, nil
}

// pearlToSign builds Pearl's to_spend / to_sign pair and returns to_sign.
func pearlToSign(msg string, pkScript []byte) (*wire.MsgTx, error) {
	toSpend := wire.NewMsgTx(2)
	toSpend.AddTxIn(wire.NewTxIn(wire.NewOutPoint(&chainhash.Hash{}, wire.MaxPrevOutIndex), nil, nil))
	toSpend.AddTxOut(wire.NewTxOut(0, pkScript))

	messageScript, err := btc_txscript.NullDataScript([]byte(msg))
	if err != nil {
		return nil, err
	}
	toSpendHash := toSpend.TxHash()
	toSign := wire.NewMsgTx(2)
	toSign.AddTxIn(wire.NewTxIn(wire.NewOutPoint(&toSpendHash, 0), nil, nil))
	toSign.AddTxOut(wire.NewTxOut(0, messageScript))
	return toSign, nil
}
