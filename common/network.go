package common

import "strings"

// NormalizeNetwork returns the canonical table key for a CSV network value:
// upper-case, "-" folded to "_" (monthly reports write ETH_LINEA, Deloitte
// reports write ETH-LINEA).
func NormalizeNetwork(network string) string {
	return strings.ToUpper(strings.ReplaceAll(strings.TrimSpace(network), "-", "_"))
}

// NetworkType resolves the signature-verification scheme for a network.
func NetworkType(network string) (string, bool) {
	t, ok := PorNetworkTypeMap[NormalizeNetwork(network)]
	return t, ok
}

// NetworkAddressType resolves the address-derivation type consumed by the
// Verify* functions. All EVM networks share ETH-style addresses; everything
// else defaults to the network name itself.
func NetworkAddressType(network string) string {
	n := NormalizeNetwork(network)
	if t, ok := networkAddressTypeExceptions[n]; ok {
		return t
	}
	if PorNetworkTypeMap[n] == EvmCoinTye {
		return "ETH"
	}
	return n
}

// MsgHeaderForNetwork derives the message signature header from the
// network's verification scheme. Unknown networks fall through to the OKX
// header, matching the Deloitte flow's default-ECDSA handling.
func MsgHeaderForNetwork(network string) string {
	n := NormalizeNetwork(network)
	switch PorNetworkTypeMap[n] {
	case EvmCoinTye, BethCoinType:
		return EthMessageSignatureHeader
	case TrxCoinType:
		return TronMessageSignatureHeader
	case UTXOCoinType:
		if h, ok := utxoMsgHeaderExceptions[n]; ok {
			return h
		}
		return BtcMessageSignatureHeader
	default:
		return OKXMessageSignatureHeader
	}
}

var (
	// PorNetworkTypeMap maps a normalized network name to its verification
	// scheme. Keys cover the monthly report Network column, the Deloitte
	// report network column, and chain-name keys inherited from the legacy
	// coin-keyed maps. Deloitte networks that the legacy maps never knew
	// (ARK, AVAX, CKB, ...) are intentionally absent: the Deloitte flow
	// falls back to ECDSA-with-pubkey for them.
	PorNetworkTypeMap = map[string]string{
		// UTXO
		"BTC": UTXOCoinType, "BCH": UTXOCoinType, "BCHN": UTXOCoinType,
		"BCHA": UTXOCoinType, "BSV": UTXOCoinType, "BCHSV": UTXOCoinType,
		"LTC": UTXOCoinType, "DOGE": UTXOCoinType, "DASH": UTXOCoinType,
		"BTG": UTXOCoinType, "BCD": UTXOCoinType, "DGB": UTXOCoinType,
		"QTUM": UTXOCoinType, "RVN": UTXOCoinType, "ZEC": UTXOCoinType,
		"ZCASH": UTXOCoinType, "XEC": UTXOCoinType, "ZEN": UTXOCoinType,
		"DCR": UTXOCoinType,

		// EVM
		"ETH": EvmCoinTye, "ETC": EvmCoinTye, "ETHW": EvmCoinTye, "ETHF": EvmCoinTye,
		"ARBITRUM": EvmCoinTye, "ETH_ARBITRUM": EvmCoinTye, "ARBITRUM_NOVA": EvmCoinTye,
		"OPTIMISM": EvmCoinTye, "ETH_OPTIMISM": EvmCoinTye,
		"BASE": EvmCoinTye, "ETH_BASE": EvmCoinTye, "ETH_LINEA": EvmCoinTye,
		"ZKSYNC2": EvmCoinTye, "ETH_ZKSYNC2": EvmCoinTye, "ETH_ZKSYNC": EvmCoinTye,
		"POLYGON": EvmCoinTye, "POLYGON_ZKEVM": EvmCoinTye,
		"ETH_ROBINHOOD": EvmCoinTye,
		"XLAYER": EvmCoinTye, "OKB_X1": EvmCoinTye, "OKB_X1_TEST": EvmCoinTye,
		"AVAXC": EvmCoinTye, "BSC": EvmCoinTye, "BNB_BSC": EvmCoinTye, "OP_BNB": EvmCoinTye,
		"CFX_EVM": EvmCoinTye, "STORY": EvmCoinTye, "IP_STORY": EvmCoinTye,
		"MERLIN_BTC": EvmCoinTye, "MERLIN": EvmCoinTye,
		"CHZ": EvmCoinTye, "CHILIZ": EvmCoinTye,
		"RONIN": EvmCoinTye, "RON": EvmCoinTye,
		"KAIA": EvmCoinTye, "KLAY": EvmCoinTye,
		"ACE": EvmCoinTye, "ACE_MAIN": EvmCoinTye,
		"MOONBEAM": EvmCoinTye, "MOVR": EvmCoinTye,
		"METIS": EvmCoinTye, "METIS_MAINNET": EvmCoinTye,
		"FLR": EvmCoinTye, "FLARE": EvmCoinTye,
		"SCROLL": EvmCoinTye, "THETA": EvmCoinTye, "TFUEL": EvmCoinTye,
		"SONIC": EvmCoinTye, "S_SONIC": EvmCoinTye,
		"CELO": EvmCoinTye, "G": EvmCoinTye, "G_GRAVITY": EvmCoinTye, "GLRM": EvmCoinTye,
		"SOPHON": EvmCoinTye, "OASYS": EvmCoinTye, "OAS": EvmCoinTye, "ZETA": EvmCoinTye,
		"POLY_OPERA": EvmCoinTye, "FTM_OPERA": EvmCoinTye, "CRONOS": EvmCoinTye,
		"MXC_MXC": EvmCoinTye, "FIL_EVM": EvmCoinTye, "FEVM": EvmCoinTye,
		"XPL": EvmCoinTye, "SEI_EVM": EvmCoinTye, "MONAD": EvmCoinTye,
		"HYPEREVM": EvmCoinTye, "GUSDT_STABLE": EvmCoinTye, "ETH_KATANA": EvmCoinTye,
		"FRAX_EVM": EvmCoinTye, "ABSTRACT": EvmCoinTye, "CTXC": EvmCoinTye,
		"ETH_WORLD": EvmCoinTye, "SENTIENT": EvmCoinTye, "ETH_MEGA": EvmCoinTye,
		"PLS": EvmCoinTye, "HT_HECO": EvmCoinTye, "ETH_BLAST": EvmCoinTye,
		"APE_APECHAIN": EvmCoinTye, "MANTA2": EvmCoinTye, "TAIKO": EvmCoinTye,
		"KCC": EvmCoinTye, "ISLM": EvmCoinTye, "HOO": EvmCoinTye,
		"FITFI_STEP": EvmCoinTye, "FAN": EvmCoinTye, "ASTR_ASTAREVM": EvmCoinTye,
		"YOU": EvmCoinTye, "XETA": EvmCoinTye, "WTC": EvmCoinTye,
		"WEMIX3": EvmCoinTye, "TRUE": EvmCoinTye, "OMN": EvmCoinTye,
		"LINKEYE": EvmCoinTye, "INT_NEW": EvmCoinTye, "FSN": EvmCoinTye,
		"EVMOS": EvmCoinTye, "CMT": EvmCoinTye, "AAC_NEW": EvmCoinTye,
		"ETH_MODE": EvmCoinTye, "STC_EVM": EvmCoinTye, "OVER": EvmCoinTye,
		"DYM": EvmCoinTye, "DYM_TEST": EvmCoinTye, "WIN": EvmCoinTye,
		"VNT": EvmCoinTye, "LTK": EvmCoinTye, "IOTX": EvmCoinTye, "LAT": EvmCoinTye,
		"ETH_LSK": EvmCoinTye, "CORE": EvmCoinTye, "BERA": EvmCoinTye,
		"UNICHAIN": EvmCoinTye, "ETH_UNICHAIN": EvmCoinTye,

		// BETH
		"BETH": BethCoinType,

		// TRX
		"TRX": TrxCoinType, "TRON": TrxCoinType,

		// ECDSA
		"FIL": EcdsaCoinType, "CFX": EcdsaCoinType,
		"ELF": EcdsaCoinType, "AELF": EcdsaCoinType,
		"LUNC": EcdsaCoinType, "LUNA_TERRA": EcdsaCoinType, "TERRA": EcdsaCoinType,
		"OKT": EcdsaCoinType, "OKC": EcdsaCoinType,
		"XRP": EcdsaCoinType, "RIPPLE": EcdsaCoinType, "NULS": EcdsaCoinType,
		"STX": EcdsaCoinType, "TIA": EcdsaCoinType, "ATOM": EcdsaCoinType,
		"CRO": EcdsaCoinType, "DORA": EcdsaCoinType, "DYDX": EcdsaCoinType,
		"INJ": EcdsaCoinType, "ICX": EcdsaCoinType,
		"ONE": EcdsaCoinType, "HARMONY_ONE": EcdsaCoinType, "FLOW": EcdsaCoinType,

		// ED25519
		"SOL": Ed25519CoinType, "FOGO": Ed25519CoinType,
		"APTOS": Ed25519CoinType, "SUI": Ed25519CoinType,
		"TON": Ed25519CoinType, "TONCOIN_NEW": Ed25519CoinType,
		"VENOM": Ed25519CoinType,
		"DOT": Ed25519CoinType, "ASSET_HUB": Ed25519CoinType,
		"ASSET_HUB_KSM": Ed25519CoinType, "KSM": Ed25519CoinType,
		"ENJIN": Ed25519CoinType, "ENJ": Ed25519CoinType,
		"PHA_KHA": Ed25519CoinType, "PHALA": Ed25519CoinType,
		"CLV": Ed25519CoinType, "AVAIL": Ed25519CoinType, "SDN": Ed25519CoinType,
		"CFG_DOT": Ed25519CoinType, "EFI_DOT": Ed25519CoinType, "KARU": Ed25519CoinType,
		"XLM": Ed25519CoinType, "PI": Ed25519CoinType, "STELLAR": Ed25519CoinType,
		"ADA": Ed25519CoinType, "NEAR": Ed25519CoinType,
		"HBAR": Ed25519CoinType, "HEDERA": Ed25519CoinType,
		"ICP": Ed25519CoinType, "ACA": Ed25519CoinType, "CSPR": Ed25519CoinType,
		"SC": Ed25519CoinType, "IOTA": Ed25519CoinType,
		"IOTA_STARDUST": Ed25519CoinType, "IOTA_NEW": Ed25519CoinType,
		"TEZOS": Ed25519CoinType, "XTZ": Ed25519CoinType,
		"KDA": Ed25519CoinType, "EGLD": Ed25519CoinType,
		"ALGO": Ed25519CoinType, "ASTR": Ed25519CoinType, "IOST": Ed25519CoinType,
		"ERD": Ed25519CoinType, "CC_CANTON": Ed25519CoinType,
		"KLY": Ed25519CoinType, "POKT": Ed25519CoinType, "LSK": Ed25519CoinType,

		// STARK
		"STARKNET": StarkCoinType, "STARKNET_ETH": StarkCoinType, "STARK": StarkCoinType,

		// EOS
		"EOS": EOSCoinType, "A": EOSCoinType, "VAULTA": EOSCoinType,
		"WAX": EOSCoinType, "WAXP": EOSCoinType,

		// PEARL
		"PRL": PearlCoinType, "PEARL": PearlCoinType, "PEARL_MAINNET": PearlCoinType,
	}

	// utxoMsgHeaderExceptions lists UTXO networks whose message header is not
	// the Bitcoin one.
	utxoMsgHeaderExceptions = map[string]string{
		"LTC": LtcMessageSignatureHeader, "DOGE": DogeMessageSignatureHeader,
		"DASH": DashMessageSignatureHeader, "BTG": BtgMessageSignatureHeader,
		"BCD": BcdMessageSignatureHeader, "DGB": DgbMessageSignatureHeader,
		"QTUM": QtumMessageSignatureHeader, "RVN": RvnMessageSignatureHeader,
		"ZEC": ZecMessageSignatureHeader, "ZCASH": ZecMessageSignatureHeader,
	}

	// networkAddressTypeExceptions lists networks whose address-derivation
	// type differs from both the network name and the EVM default.
	networkAddressTypeExceptions = map[string]string{
		"BCHN": "BCH", "XEC": "BCH", "ZCASH": "ZEC",
		"BSV": "BTC", "BCD": "BTC", "BCHSV": "BTC",
		"FIL_EVM": "FIL", "FEVM": "FIL", "LAT": "LAT",
		"OKC": "ETH", "OKT": "ETH",
		"AELF": "ELF",
		"TRON": "TRX",
		"TONCOIN_NEW": "TON",
		"ASSET_HUB": "DOT", "ASSET_HUB_KSM": "KSM",
		"PHA_KHA": "PHA", "PHALA": "PHA",
		"CFG_DOT": "CFG", "EFI_DOT": "EFI",
		"HEDERA": "HBAR",
		"LUNA_TERRA": "TERRA",
		"HARMONY_ONE": "ONE",
		"A": "EOS", "VAULTA": "EOS", "WAX": "EOS", "WAXP": "EOS",
		"FOGO": "SOL",
	}
)
