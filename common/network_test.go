package common

import "testing"

func TestNetworkTypeSamples(t *testing.T) {
	cases := []struct{ network, want string }{
		{"ETH", EvmCoinTye}, {"eth-arbitrum", EvmCoinTye}, {"ARBITRUM", EvmCoinTye},
		{"XLAYER", EvmCoinTye}, {"FEVM", EvmCoinTye}, {"ETC", EvmCoinTye},
		{"ETH-ROBINHOOD", EvmCoinTye},
		{"BTC", UTXOCoinType}, {"BCHN", UTXOCoinType}, {"ZCASH", UTXOCoinType},
		{"TRON", TrxCoinType}, {"TRX", TrxCoinType},
		{"OKC", EcdsaCoinType}, {"OKT", EcdsaCoinType}, {"AELF", EcdsaCoinType},
		{"XRP", EcdsaCoinType}, {"FIL", EcdsaCoinType},
		{"APTOS", Ed25519CoinType}, {"SOL", Ed25519CoinType}, {"TON", Ed25519CoinType},
		{"ASSET_HUB", Ed25519CoinType}, {"ASSET-HUB-KSM", Ed25519CoinType},
		{"LSK", Ed25519CoinType}, {"ADA", Ed25519CoinType},
		{"BETH", BethCoinType},
		{"STARKNET", StarkCoinType}, {"STARK", StarkCoinType},
		{"Vaulta", EOSCoinType}, {"WAXP", EOSCoinType},
		{"PRL", PearlCoinType}, {"PEARL", PearlCoinType}, {"PEARL-MAINNET", PearlCoinType},
	}
	for _, c := range cases {
		got, ok := NetworkType(c.network)
		if !ok || got != c.want {
			t.Errorf("NetworkType(%q) = %q, %v; want %q", c.network, got, ok, c.want)
		}
	}
	if _, ok := NetworkType("NO_SUCH_NETWORK"); ok {
		t.Error("NetworkType(NO_SUCH_NETWORK) should not be found")
	}
}

func TestNetworkAddressTypeSamples(t *testing.T) {
	cases := []struct{ network, want string }{
		{"ETH", "ETH"}, {"ARBITRUM", "ETH"}, {"XLAYER", "ETH"},
		{"FEVM", "FIL"}, {"FIL-EVM", "FIL"}, {"LAT", "LAT"},
		{"OKC", "ETH"}, {"OKT", "ETH"}, {"AELF", "ELF"},
		{"ASSET_HUB", "DOT"}, {"ASSET-HUB-KSM", "KSM"}, {"TONCOIN-NEW", "TON"},
		{"SOL", "SOL"}, {"ADA", "ADA"}, {"APTOS", "APTOS"},
		{"BCHN", "BCH"}, {"XEC", "BCH"}, {"ZCASH", "ZEC"}, {"BCHSV", "BTC"},
		{"HARMONY-ONE", "ONE"}, {"LUNA_TERRA", "TERRA"}, {"HEDERA", "HBAR"},
	}
	for _, c := range cases {
		if got := NetworkAddressType(c.network); got != c.want {
			t.Errorf("NetworkAddressType(%q) = %q; want %q", c.network, got, c.want)
		}
	}
}

func TestMsgHeaderForNetworkSamples(t *testing.T) {
	cases := []struct{ network, want string }{
		{"ETH", EthMessageSignatureHeader}, {"ARBITRUM", EthMessageSignatureHeader},
		{"BETH", EthMessageSignatureHeader},
		{"TRON", TronMessageSignatureHeader}, {"TRX", TronMessageSignatureHeader},
		{"BTC", BtcMessageSignatureHeader}, {"LTC", LtcMessageSignatureHeader},
		{"DOGE", DogeMessageSignatureHeader}, {"ZCASH", ZecMessageSignatureHeader},
		{"XEC", BtcMessageSignatureHeader}, {"ZEN", BtcMessageSignatureHeader},
		{"SOL", OKXMessageSignatureHeader}, {"OKC", OKXMessageSignatureHeader},
		{"LSK", OKXMessageSignatureHeader},
		// 未知 network：德勤流程走 ECDSA 兜底，消息头也要落到 OKX
		{"NO_SUCH_NETWORK", OKXMessageSignatureHeader},
	}
	for _, c := range cases {
		if got := MsgHeaderForNetwork(c.network); got != c.want {
			t.Errorf("MsgHeaderForNetwork(%q) = %q; want %q", c.network, got, c.want)
		}
	}
}
