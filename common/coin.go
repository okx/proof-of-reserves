package common

const (
	UTXOCoinType    = "UTXO"
	EvmCoinTye      = "EVM"
	EcdsaCoinType   = "ECDSA"
	Ed25519CoinType = "ED25519"
	StarkCoinType   = "STARK"

	TrxCoinType  = "TRX"
	BethCoinType = "BETH"
	AlgoCoinType = "ALGO"
	EOSCoinType  = "EOS"
	// PearlCoinType: Taproot script-path 2-of-2 proofs (two signatures plus tapscript).
	PearlCoinType = "PEARL"

	BtcMessageSignatureHeader  = "Bitcoin Signed Message:\n"
	LtcMessageSignatureHeader  = "Litecoin Signed Message:\n"
	DogeMessageSignatureHeader = "Dogecoin Signed Message:\n"
	DashMessageSignatureHeader = "DarkCoin Signed Message:\n"
	BtgMessageSignatureHeader  = "Bitcoin Gold Signed Message:\n"
	BcdMessageSignatureHeader  = "Bitcoindiamond Signed Message:\n"
	DgbMessageSignatureHeader  = "DigiByte Signed Message:\n"
	QtumMessageSignatureHeader = "Qtum Signed Message:\n"
	RvnMessageSignatureHeader  = "Raven Signed Message:\n"
	ZecMessageSignatureHeader  = "Zcash Signed Message:\n"

	EthMessageSignatureHeader    = "\x19Ethereum Signed Message:\n32"
	TronMessageSignatureHeader   = "\x19TRON Signed Message:\n32"
	TronMessageV2SignatureHeader = "\x19TRON Signed Message:\n"

	OKXMessageSignatureHeader = "OKX Signed Message:\n"
)

var (
	PorCoinUnitMap = map[string]string{
		// UTXO
		"BTC":       "BTC",
		"USDT-OMNI": "USDT",
		"BCHN":      "BCHN",
		"BCHA":      "BCHA",
		"BSV":       "BSV",
		"LTC":       "LTC",
		"DOGE":      "DOGE",
		"DASH":      "DASH",
		"BTG":       "BTG",
		"BCD":       "BCD",
		"DGB":       "DGB",
		"QTUM":      "QTUM",
		"RVN":       "RVN",
		"ZEC":       "ZEC",
		"PRL":       "PRL",

		// ETH
		"ETH":                  "ETH",
		"ETH-ARBITRUM":         "ETH",
		"ETH-OPTIMISM":         "ETH",
		"USDT":                 "USDT",
		"USDT-ERC20":           "USDT",
		"USDG-ERC20":           "USDG",
		"USDG-XLAYER":          "USDG",
		"USDG-ROBINHOOD-D168":  "USDG",
		"USDT-POLY":            "USDT",
		"USDT-AVAXC":           "USDT",
		"USDT-ARBITRUM":        "USDT",
		"USDT-OPTIMISM":        "USDT",
		"USDT-OKC20":           "USDT",
		"USDC":                 "USDC",
		"POLY-USDC":            "USDC",
		"USDC-AVAXC":           "USDC",
		"USDC-ARBITRUM":        "USDC",
		"USDC-OPTIMISM":        "USDC",
		"ETC":                  "ETC",
		"OKB-OKC20":            "OKB",
		"LTCK-OKC20":           "LTC",
		"FILK-OKC20":           "FIL",
		"USDC-OKC20":           "USDC",
		"SHIBK-KIP20":          "SHIB",
		"DOTK-OKC20":           "DOT",
		"ETCK-KIP20":           "ETC",
		"XRPK-KIP20":           "RIPPLE",
		"UNIK-OKC20":           "UNI",
		"BCHK-KIP20":           "BCH",
		"BABYDOGE-KIP20":       "BABYDOGE",
		"LINKK-OKC20":          "LINK",
		"TRXK-KIP20":           "TRX",
		"BABYDOGE-BSC":         "BABYDOGE",
		"SHIB":                 "SHIB",
		"UNI":                  "UNI",
		"LINK":                 "LINK",
		"ETHW":                 "ETHW",
		"BLUR":                 "BLUR",
		"MATIC":                "MATIC",
		"PEOPLE":               "PEOPLE",
		"OKT":                  "OKT",
		"OKB":                  "OKB",
		"OPTIMISM":             "OPTIMISM",
		"ETH-LINEA":            "ETH",
		"BASE":                 "ETH",
		"ETH-BASE":             "ETH",
		"ETH-ZKSYNC2":          "ETH",
		"OKB-X1":               "OKB",
		"OKB-X1-ETH":           "ETH",
		"OKB-X1-USDT":          "USDT",
		"OKB-X1-USDC":          "USDC",
		"POLY-USDC-3359":       "USDC",
		"USDC-OPTIMISM-FF85":   "USDC",
		"USDC-ARBITRUM-NATIVE": "USDC",
		"USDC-BASE":            "USDC",
		"ZKSYNC2":              "ETH",
		"ETHK-OKC20":           "ETH",
		"STARKNET-ETH":         "ETH",

		// BETH
		"BETH": "BETH",

		// TRX
		"USDT-TRC20": "USDT",
		"USDC-TRC":   "USDC",
		"TRX":        "TRX",

		// ECDSA
		"FIL":  "FIL",
		"CFX":  "CFX",
		"ELF":  "ELF",
		"LUNC": "LUNC",

		// ED25519
		"SOL":         "SOL",
		"USDC-SPL":    "USDC",
		"USDT-SPL":    "USDT",
		"USDG":        "USDG",
		"APTOS":       "APTOS",
		"APTOS-FA":    "APTOS",
		"USDT-APTOS":  "USDT",
		"USDC-APTOS":  "USDC",
		"SUI":         "SUI",
		"USDC-SUI":    "USDC",
		"TONCOIN-NEW": "TONCOIN-NEW",
		"USDT-TON":    "USDT",
		"DOT":         "DOT",
		"ASSET-HUB":   "DOT",
		"XLM":         "XLM",
		"PI":          "PI",
		"ADA":         "ADA",
		"NEAR":        "NEAR",
		"HBAR":        "HBAR",

		// EOS
		"EOS":    "EOS",
		"A":      "A",
		"RIPPLE": "RIPPLE",

		// ALGO
		"USDT-ALGO": "USDT",

		// FEVM
		"FIL-EVM": "FIL",
	}

	PorCoinBaseUnitPrecisionMap = map[string]int{
		// UTXO
		"BTC":       8,
		"USDT-OMNI": 6,
		"BCHN":      8,
		"BCHA":      8,
		"BSV":       8,
		"LTC":       8,
		"DOGE":      8,
		"DASH":      8,
		"BTG":       8,
		"BCD":       8,
		"DGB":       8,
		"QTUM":      8,
		"RVN":       8,
		"PRL":       8,
		"ZEC":       8,

		// ETH
		"ETH":                  18,
		"ETH-ARBITRUM":         18,
		"ETH-OPTIMISM":         18,
		"USDT":                 6,
		"USDT-ERC20":           6,
		"USDG-ROBINHOOD-D168":  6,
		"USDT-POLY":            6,
		"USDT-AVAXC":           6,
		"USDT-ARBITRUM":        6,
		"USDT-OPTIMISM":        6,
		"USDT-OKC20":           18,
		"USDC":                 6,
		"POLY-USDC":            6,
		"USDC-AVAXC":           6,
		"USDC-ARBITRUM":        6,
		"USDC-OPTIMISM":        6,
		"ETC":                  18,
		"OKB-OKC20":            18,
		"LTCK-OKC20":           18,
		"FILK-OKC20":           18,
		"USDC-OKC20":           18,
		"SHIBK-KIP20":          18,
		"DOTK-OKC20":           18,
		"ETCK-KIP20":           18,
		"XRPK-KIP20":           18,
		"UNIK-OKC20":           18,
		"BCHK-KIP20":           18,
		"BABYDOGE-KIP20":       18,
		"LINKK-OKC20":          18,
		"TRXK-KIP20":           18,
		"BABYDOGE-BSC":         18,
		"SHIB":                 18,
		"UNI":                  18,
		"LINK":                 18,
		"ETHW":                 18,
		"BLUR":                 18,
		"MATIC":                18,
		"PEOPLE":               18,
		"OKT":                  18,
		"OKB":                  18,
		"OPTIMISM":             18,
		"ETH-LINEA":            18,
		"BASE":                 18,
		"ETH-BASE":             18,
		"ETH-ZKSYNC2":          18,
		"OKB-X1":               18,
		"OKB-X1-ETH":           18,
		"OKB-X1-USDT":          6,
		"OKB-X1-USDC":          6,
		"POLY-USDC-3359":       6,
		"USDC-OPTIMISM-FF85":   6,
		"USDC-ARBITRUM-NATIVE": 6,
		"USDC-BASE":            6,
		"ZKSYNC2":              18,
		"STARKNET-ETH":         18,

		// BETH
		"BETH": 18,

		// TRX
		"USDT-TRC20": 6,
		"USDC-TRC":   6,
		"TRX":        6,

		// ECDSA
		"FIL":  18,
		"CFX":  18,
		"ELF":  8,
		"LUNC": 6,

		// ED25519
		"SOL":         9,
		"USDC-SPL":    6,
		"USDT-SPL":    6,
		"USDG":        6,
		"APTOS":       8,
		"APTOS-FA":    8,
		"USDT-APTOS":  6,
		"USDC-APTOS":  6,
		"SUI":         9,
		"USDC-SUI":    6,
		"TONCOIN-NEW": 9,
		"USDT-TON":    6,
		"DOT":         10,
		"ASSET-HUB":   10,

		// EOS
		"EOS":    4,
		"A":      4,
		"RIPPLE": 6,

		// ALGO
		"USDT-ALGO": 6,

		// FEVM
		"FIL-EVM": 18,
	}

	CheckBalanceCoinBlackList = map[string]bool{
		"DASH": true,
		"DOGE": true,
		"BCHN": true,

		"TRX":      true,
		"USDC-TRC": true,
		"BETH":     true,
		"ETC":      true,

		"FIL":  true,
		"CFX":  true,
		"ELF":  true,
		"LUNC": true,

		"SOL":         true,
		"USDC-SPL":    true,
		"APTOS":       true,
		"TONCOIN-NEW": true,
		"DOT":         true,
		"NEAR":        true,
		"HBAR":        true,

		"RIPPLE": true,

		"USDT-ALGO": true,
		"FIL-EVM":   true,
		"ETH-LINEA": true,
		"BASE":      true,

		"OKB-X1":      true,
		"OKB-X1-ETH":  true,
		"OKB-X1-USDT": true,
		"OKB-X1-USDC": true,
	}

	VerifyAddressCoinBlackList = map[string]bool{
		"EOS":       true,
		"A":         true,
		"RIPPLE":    true,
		"USDT-ALGO": true,
	}
)

func IsCheckBalanceBannedCoin(coin string) bool {
	return CheckBalanceCoinBlackList[coin]
}

func IsVerifyAddressBannedCoin(coin string) bool {
	return VerifyAddressCoinBlackList[coin]
}
