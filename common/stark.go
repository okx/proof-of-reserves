package common

import (
	"fmt"
	"sync"

	"github.com/okx/go-wallet-sdk/coins/starknet"
)

var starknetMutex sync.Mutex

func VerifyStarknetEIP712(accountAddress, msg, publicKey, sig string) bool {
	starknetMutex.Lock()
	defer starknetMutex.Unlock()

	const EIP712_TEMPLATE = `{
    "accountAddress": "%s",
    "typedData": {
        "types": {
            "StarkNetDomain": [
                {
                    "name": "name",
                    "type": "felt"
                },
                {
                    "name": "version",
                    "type": "felt"
                },
                {
                    "name": "chainId",
                    "type": "felt"
                }
            ],
            "Message": [
                {
                    "name": "contents",
                    "type": "felt"
                }
            ]
        },
        "primaryType": "Message",
        "domain": {
            "name": "OKX POR MESSAGE",
            "version": "1",
            "chainId": "0x534e5f4d41494e"
        },
        "message": {
            "contents": "%s"
        }
    }
}`

	hash, err := starknet.GetMessageHashWithJson(fmt.Sprintf(EIP712_TEMPLATE, accountAddress, msg))
	if err != nil || len(hash) == 0 {
		return false
	}

	curve := starknet.SC()

	sigR := starknet.HexToBig(sig[:64])
	sigS := starknet.HexToBig(sig[64:])
	pubX, pubY := curve.XToPubKey(publicKey)

	return curve.Verify(starknet.HexToBig(hash), sigR, sigS, pubX, pubY)
}
