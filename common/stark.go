package common

import "github.com/okx/go-wallet-sdk/coins/starknet"

var curve = starknet.SC()

func VerifyMessage(hash, publicKey, sig string) bool {
	sigR := starknet.HexToBig(sig[:64])
	sigS := starknet.HexToBig(sig[64:])
	pubX, pubY := curve.XToPubKey(publicKey)

	return curve.Verify(starknet.HexToBig(hash), sigR, sigS, pubX, pubY)
}
