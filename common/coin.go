package common

const (
	UTXOCoinType    = "UTXO"
	EvmCoinTye      = "EVM"
	EcdsaCoinType   = "ECDSA"
	Ed25519CoinType = "ED25519"

	TrxCoinType  = "TRX"
	BethCoinType = "BETH"
	AlgoCoinType = "ALGO"

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
	PorCoinAddressTypeMap = map[string]string{
		// UTXO
		"BTC":       "BTC",
		"USDT-OMNI": "BTC",
		"BCHN":      "BCH",
		"BCHA":      "BCHA",
		"BSV":       "BTC",
		"LTC":       "LTC",
		"DOGE":      "DOGE",
		"DASH":      "DASH",
		"BTG":       "BTG",
		"BCD":       "BTC",
		"DGB":       "DGB",
		"QTUM":      "QTUM",
		"RVN":       "RVN",
		"ZEC":       "ZEC",

		// ETH
		"ETH":                  "ETH",
		"ETH-ARBITRUM":         "ETH",
		"ETH-OPTIMISM":         "ETH",
		"USDT-ERC20":           "ETH",
		"USDT-POLY":            "ETH",
		"USDT-AVAXC":           "ETH",
		"USDT-ARBITRUM":        "ETH",
		"USDT-OPTIMISM":        "ETH",
		"USDT-OKC20":           "ETH",
		"USDC":                 "ETH",
		"POLY-USDC":            "ETH",
		"USDC-AVAXC":           "ETH",
		"USDC-ARBITRUM":        "ETH",
		"USDC-OPTIMISM":        "ETH",
		"ETC":                  "ETH",
		"OKB-OKC20":            "ETH",
		"LTCK-OKC20":           "ETH",
		"FILK-OKC20":           "ETH",
		"USDC-OKC20":           "ETH",
		"SHIBK-KIP20":          "ETH",
		"DOTK-OKC20":           "ETH",
		"ETCK-KIP20":           "ETH",
		"XRPK-KIP20":           "ETH",
		"UNIK-OKC20":           "ETH",
		"BCHK-KIP20":           "ETH",
		"BABYDOGE-KIP20":       "ETH",
		"LINKK-OKC20":          "ETH",
		"TRXK-KIP20":           "ETH",
		"BABYDOGE-BSC":         "ETH",
		"SHIB":                 "ETH",
		"UNI":                  "ETH",
		"LINK":                 "ETH",
		"ETHW":                 "ETH",
		"BLUR":                 "ETH",
		"MATIC":                "ETH",
		"PEOPLE":               "ETH",
		"OKT":                  "ETH",
		"OKB":                  "ETH",
		"OPTIMISM":             "ETH",
		"ETH-LINEA":            "ETH",
		"BASE":                 "ETH",
		"OKB-X1":               "ETH",
		"OKB-X1-ETH":           "ETH",
		"OKB-X1-USDT":          "ETH",
		"OKB-X1-USDC":          "ETH",
		"POLY-USDC-3359":       "ETH",
		"USDC-OPTIMISM-FF85":   "ETH",
		"USDC-ARBITRUM-NATIVE": "ETH",
		"USDC-BASE":            "ETH",

		// BETH
		"BETH": "BETH",

		// TRX
		"USDT-TRC20": "TRX",
		"USDC-TRC":   "TRX",
		"TRX":        "TRX",

		// ECDSA
		"FIL":  "FIL",
		"CFX":  "CFX",
		"ELF":  "ELF",
		"LUNC": "LUNC",

		// ED25519
		"SOL":         "SOL",
		"USDC-SPL":    "SOL",
		"APTOS":       "APTOS",
		"TONCOIN-NEW": "TON",
		"DOT":         "DOT",

		// ALGO
		"USDT-ALGO": "ALGO",

		// FEVM
		"FIL-EVM": "FIL",
	}

	PorCoinTypeMap = map[string]string{
		// UTXO
		"BTC":       UTXOCoinType,
		"USDT-OMNI": UTXOCoinType,
		"BCHN":      UTXOCoinType,
		"BCHA":      UTXOCoinType,
		"BSV":       UTXOCoinType,
		"LTC":       UTXOCoinType,
		"DOGE":      UTXOCoinType,
		"DASH":      UTXOCoinType,
		"BTG":       UTXOCoinType,
		"BCD":       UTXOCoinType,
		"DGB":       UTXOCoinType,
		"QTUM":      UTXOCoinType,
		"RVN":       UTXOCoinType,
		"ZEC":       UTXOCoinType,

		// ETH
		"ETH":                  EvmCoinTye,
		"ETH-ARBITRUM":         EvmCoinTye,
		"ETH-OPTIMISM":         EvmCoinTye,
		"USDT-ERC20":           EvmCoinTye,
		"USDT-POLY":            EvmCoinTye,
		"USDT-AVAXC":           EvmCoinTye,
		"USDT-ARBITRUM":        EvmCoinTye,
		"USDT-OPTIMISM":        EvmCoinTye,
		"USDC":                 EvmCoinTye,
		"POLY-USDC":            EvmCoinTye,
		"USDC-AVAXC":           EvmCoinTye,
		"USDC-ARBITRUM":        EvmCoinTye,
		"USDC-OPTIMISM":        EvmCoinTye,
		"ETC":                  EvmCoinTye,
		"BABYDOGE-BSC":         EvmCoinTye,
		"SHIB":                 EvmCoinTye,
		"UNI":                  EvmCoinTye,
		"LINK":                 EvmCoinTye,
		"ETHW":                 EvmCoinTye,
		"BLUR":                 EvmCoinTye,
		"MATIC":                EvmCoinTye,
		"PEOPLE":               EvmCoinTye,
		"OKB":                  EvmCoinTye,
		"OPTIMISM":             EvmCoinTye,
		"FIL-EVM":              EvmCoinTye,
		"ETH-LINEA":            EvmCoinTye,
		"BASE":                 EvmCoinTye,
		"OKB-X1":               EvmCoinTye,
		"OKB-X1-ETH":           EvmCoinTye,
		"OKB-X1-USDT":          EvmCoinTye,
		"OKB-X1-USDC":          EvmCoinTye,
		"POLY-USDC-3359":       EvmCoinTye,
		"USDC-OPTIMISM-FF85":   EvmCoinTye,
		"USDC-ARBITRUM-NATIVE": EvmCoinTye,
		"USDC-BASE":            EvmCoinTye,

		// BETH
		"BETH": BethCoinType,

		// TRX
		"USDT-TRC20": TrxCoinType,
		"USDC-TRC":   TrxCoinType,
		"TRX":        TrxCoinType,

		// ECDSA
		"FIL":            EcdsaCoinType,
		"CFX":            EcdsaCoinType,
		"ELF":            EcdsaCoinType,
		"LUNC":           EcdsaCoinType,
		"OKB-OKC20":      EcdsaCoinType,
		"LTCK-OKC20":     EcdsaCoinType,
		"FILK-OKC20":     EcdsaCoinType,
		"USDC-OKC20":     EcdsaCoinType,
		"SHIBK-KIP20":    EcdsaCoinType,
		"DOTK-OKC20":     EcdsaCoinType,
		"ETCK-KIP20":     EcdsaCoinType,
		"XRPK-KIP20":     EcdsaCoinType,
		"UNIK-OKC20":     EcdsaCoinType,
		"BCHK-KIP20":     EcdsaCoinType,
		"BABYDOGE-KIP20": EcdsaCoinType,
		"LINKK-OKC20":    EcdsaCoinType,
		"TRXK-KIP20":     EcdsaCoinType,
		"OKT":            EcdsaCoinType,
		"USDT-OKC20":     EcdsaCoinType,

		// ED25519
		"SOL":         Ed25519CoinType,
		"USDC-SPL":    Ed25519CoinType,
		"APTOS":       Ed25519CoinType,
		"TONCOIN-NEW": Ed25519CoinType,
		"DOT":         Ed25519CoinType,

		// ALGO
		"USDT-ALGO": AlgoCoinType,
	}

	PorCoinMessageSignatureHeaderMap = map[string]string{
		// UTXO
		"BTC":       BtcMessageSignatureHeader,
		"USDT-OMNI": BtcMessageSignatureHeader,
		"BCHN":      BtcMessageSignatureHeader,
		"BCHA":      BtcMessageSignatureHeader,
		"BSV":       BtcMessageSignatureHeader,
		"LTC":       LtcMessageSignatureHeader,
		"DOGE":      DogeMessageSignatureHeader,
		"DASH":      DashMessageSignatureHeader,
		"BTG":       BtgMessageSignatureHeader,
		"BCD":       BcdMessageSignatureHeader,
		"DGB":       DgbMessageSignatureHeader,
		"QTUM":      QtumMessageSignatureHeader,
		"RVN":       RvnMessageSignatureHeader,
		"ZEC":       ZecMessageSignatureHeader,

		// ETH
		"ETH":                  EthMessageSignatureHeader,
		"ETH-ARBITRUM":         EthMessageSignatureHeader,
		"ETH-OPTIMISM":         EthMessageSignatureHeader,
		"USDT-ERC20":           EthMessageSignatureHeader,
		"USDT-POLY":            EthMessageSignatureHeader,
		"USDT-AVAXC":           EthMessageSignatureHeader,
		"USDT-ARBITRUM":        EthMessageSignatureHeader,
		"USDT-OPTIMISM":        EthMessageSignatureHeader,
		"USDC":                 EthMessageSignatureHeader,
		"POLY-USDC":            EthMessageSignatureHeader,
		"USDC-AVAXC":           EthMessageSignatureHeader,
		"USDC-ARBITRUM":        EthMessageSignatureHeader,
		"USDC-OPTIMISM":        EthMessageSignatureHeader,
		"ETC":                  EthMessageSignatureHeader,
		"BABYDOGE-BSC":         EthMessageSignatureHeader,
		"SHIB":                 EthMessageSignatureHeader,
		"UNI":                  EthMessageSignatureHeader,
		"LINK":                 EthMessageSignatureHeader,
		"ETHW":                 EthMessageSignatureHeader,
		"BLUR":                 EthMessageSignatureHeader,
		"MATIC":                EthMessageSignatureHeader,
		"PEOPLE":               EthMessageSignatureHeader,
		"OKB":                  EthMessageSignatureHeader,
		"OPTIMISM":             EthMessageSignatureHeader,
		"FIL-EVM":              EthMessageSignatureHeader,
		"ETH-LINEA":            EthMessageSignatureHeader,
		"BASE":                 EthMessageSignatureHeader,
		"OKB-X1":               EthMessageSignatureHeader,
		"OKB-X1-ETH":           EthMessageSignatureHeader,
		"OKB-X1-USDT":          EthMessageSignatureHeader,
		"OKB-X1-USDC":          EthMessageSignatureHeader,
		"POLY-USDC-3359":       EthMessageSignatureHeader,
		"USDC-OPTIMISM-FF85":   EthMessageSignatureHeader,
		"USDC-ARBITRUM-NATIVE": EthMessageSignatureHeader,
		"USDC-BASE":            EthMessageSignatureHeader,

		// BETH
		"BETH": EthMessageSignatureHeader,

		// TRX
		"USDT-TRC20": TronMessageSignatureHeader,
		"USDC-TRC":   TronMessageSignatureHeader,
		"TRX":        TronMessageSignatureHeader,

		// ECDSA
		"FIL":            OKXMessageSignatureHeader,
		"CFX":            OKXMessageSignatureHeader,
		"ELF":            OKXMessageSignatureHeader,
		"LUNC":           OKXMessageSignatureHeader,
		"OKB-OKC20":      OKXMessageSignatureHeader,
		"LTCK-OKC20":     OKXMessageSignatureHeader,
		"FILK-OKC20":     OKXMessageSignatureHeader,
		"USDC-OKC20":     OKXMessageSignatureHeader,
		"SHIBK-KIP20":    OKXMessageSignatureHeader,
		"DOTK-OKC20":     OKXMessageSignatureHeader,
		"ETCK-KIP20":     OKXMessageSignatureHeader,
		"XRPK-KIP20":     OKXMessageSignatureHeader,
		"UNIK-OKC20":     OKXMessageSignatureHeader,
		"BCHK-KIP20":     OKXMessageSignatureHeader,
		"BABYDOGE-KIP20": OKXMessageSignatureHeader,
		"LINKK-OKC20":    OKXMessageSignatureHeader,
		"TRXK-KIP20":     OKXMessageSignatureHeader,
		"OKT":            OKXMessageSignatureHeader,
		"USDT-OKC20":     OKXMessageSignatureHeader,

		// ED25519
		"SOL":         OKXMessageSignatureHeader,
		"USDC-SPL":    OKXMessageSignatureHeader,
		"APTOS":       OKXMessageSignatureHeader,
		"TONCOIN-NEW": OKXMessageSignatureHeader,
		"DOT":         OKXMessageSignatureHeader,

		// ALGO
		"USDT-ALGO": OKXMessageSignatureHeader,
	}

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

		// ETH
		"ETH":                  "ETH",
		"ETH-ARBITRUM":         "ETH",
		"ETH-OPTIMISM":         "ETH",
		"USDT":                 "USDT",
		"USDT-ERC20":           "USDT",
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
		"OKB-X1":               "OKB",
		"OKB-X1-ETH":           "ETH",
		"OKB-X1-USDT":          "USDT",
		"OKB-X1-USDC":          "USDC",
		"POLY-USDC-3359":       "USDC",
		"USDC-OPTIMISM-FF85":   "USDC",
		"USDC-ARBITRUM-NATIVE": "USDC",
		"USDC-BASE":            "USDC",

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
		"APTOS":       "APTOS",
		"TONCOIN-NEW": "TONCOIN-NEW",
		"DOT":         "DOT",

		// EOS
		"EOS":    "EOS",
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
		"ZEC":       8,

		// ETH
		"ETH":                  18,
		"ETH-ARBITRUM":         18,
		"ETH-OPTIMISM":         18,
		"USDT":                 6,
		"USDT-ERC20":           6,
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
		"OKB-X1":               18,
		"OKB-X1-ETH":           18,
		"OKB-X1-USDT":          6,
		"OKB-X1-USDC":          6,
		"POLY-USDC-3359":       6,
		"USDC-OPTIMISM-FF85":   6,
		"USDC-ARBITRUM-NATIVE": 6,
		"USDC-BASE":            6,

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
		"APTOS":       8,
		"TONCOIN-NEW": 9,
		"DOT":         10,

		// EOS
		"EOS":    4,
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

		"EOS":    true,
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
