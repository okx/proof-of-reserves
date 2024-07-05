package common

import (
	"encoding/hex"
	"testing"
)

func TestUtxoCoinSigToPubKey(t *testing.T) {
	args := []struct {
		coin string
		addr string
		msg  string
		sign string
		want string
	}{
		{
			coin: "BTC",
			addr: "1DcT5Wij5tfb3oVViF8mA8p4WrG98ahZPT",
			msg:  "I am an OKX address",
			sign: "IA1jDx3zkn4J4F6mCVU68Vm7TwNf+bCsp+hKo3LwV/Y+PlZEoNsajnAHqd/FrEmv5/VAGz7pPiWPOXjmCLRfxIM=",
			want: "024f85415b4038658f84e316cbf0dd0eed649ff778b9440113fc8ad12832d612c9",
		},
		{
			coin: "BTC",
			addr: "3Fs7C97NmvhWUZ2pSjth9YbTxMv4sk9nHi",
			msg:  "I am a OKX address",
			sign: "H2vshvcYTGrUw0XG1AundmbivdrhTWUOTqcXKhN+MqbaEfVYfGkgDhEumiJoEJFhlzuma6bBpg4pXNUHoTENOPI=",
			want: "03447bead626f13c79de937c0879b64172e5984456a47350b44e8bd23a02e6895e",
		},
		{
			coin: "BTC",
			addr: "bc1qpypsu8sytw959yu53dk48eaq9saxumwegzwd4anava9qe40k6gfqyrsxaq",
			msg:  "I am a OKX address",
			sign: "IKfykjJJSywz2g/KGXvWwE1aphUvryiBaNQYK0m4Ain7QkNhM7VjwV964DPn4dvpOGyzhQYSAqLIz1BOBXbIQcY=",
			want: "02b514e7ccc2845d3f1ca7181dacab0d1ac277616e753547922f82cc0cdfb5c691",
		},
		{
			coin: "BCH",
			addr: "393maTY7rQScy4SmYE1XSXUSgK73byhgfA",
			msg:  "hello world",
			sign: "Hz+cZI5GfSzNSvBpna20diV47/rhlQMRQTNGZd9sI4UZQaWH4ZY3KJA4IlcP5bwuicO+myA4vLdiMkj7OU+rDpg=",
			want: "03052b16e71e4413f24f8504c3b188b7edebf97b424582877e4993ef9b23d0f045",
		},
		{
			coin: "BSV",
			addr: "1Hgc1DnWHwfXFxruej4H5g3ThsCzUEwLdD",
			msg:  "hello world",
			sign: "IDtG3XPLpiKOp4PjTzCo/ng8gm4MFTTyHeh/DaPC1XYsYaj5Jr4h8dnxmwuJtNkPkH40rEfnrrO8fgZKNOIF5iM=",
			want: "0273fa0df3ffceeda23b0074d9fe83d9ee3a209fad6e4546fdec5ede39abcbb70d",
		},
		{
			coin: "DOGE",
			addr: "9yo2KJc1vUKWsRpExMfwgf6pNtV5gG3vqW",
			msg:  "hello world",
			sign: "Hz+cZI5GfSzNSvBpna20diV47/rhlQMRQTNGZd9sI4UZQaWH4ZY3KJA4IlcP5bwuicO+myA4vLdiMkj7OU+rDpg=",
			want: "026e990eb5b0b641433908d35ac75364895e59c1e1c2002c425b281059876c2504",
		},
	}
	for _, tt := range args {
		get, err := UtxoCoinSigToPubKey(tt.coin, tt.msg, tt.sign)
		if hex.EncodeToString(get) != tt.want {
			t.Errorf("coin: %s, addr: %s, msg: %s, sign: %s, get: %s, want:%s", tt.coin, tt.addr, tt.msg, tt.sign, hex.EncodeToString(get), tt.want)
		}
		if err != nil {
			t.Errorf(err.Error())
		}
	}
}
