package common

import (
	"github.com/martinboehm/bchutil"
	"github.com/martinboehm/btcutil"
	"github.com/martinboehm/btcutil/chaincfg"
	"regexp"
	"strings"
)

var (
	// MainNetParams are parser parameters for mainnet
	MainNetParams chaincfg.Params
)

// GetBTCMainNetParams BTC
func GetBTCMainNetParams() *chaincfg.Params {
	if !chaincfg.IsRegistered(&chaincfg.MainNetParams) {
		chaincfg.RegisterBitcoinParams()
	}
	return &chaincfg.MainNetParams
}

// GetBCHMainNetParams BCH
func GetBCHMainNetParams() *chaincfg.Params {
	localParams := chaincfg.MainNetParams
	localParams.Net = 0xe8f3e1e3

	// Address encoding magics
	localParams.PubKeyHashAddrID = []byte{0} // base58 prefix: 1
	localParams.ScriptHashAddrID = []byte{5} // base58 prefix: 3

	if !chaincfg.IsRegistered(&localParams) {
		err := chaincfg.Register(&localParams)
		if err != nil {
			panic(err)
		}
	}
	return &localParams
}

// GetECASHMainNetParams ECASH
func GetECASHMainNetParams() *chaincfg.Params {
	localParams := chaincfg.MainNetParams
	localParams.Net = 0xe8f3e1e3

	// Address encoding magics
	localParams.PubKeyHashAddrID = []byte{0} // base58 prefix: 1
	localParams.ScriptHashAddrID = []byte{5} // base58 prefix: 3

	if !chaincfg.IsRegistered(&localParams) {
		err := chaincfg.Register(&localParams)
		if err != nil {
			panic(err)
		}
	}
	return &localParams
}

// GetLTCMainNetParams LTC
func GetLTCMainNetParams() *chaincfg.Params {
	localParams := chaincfg.MainNetParams
	localParams.Net = 0xdbb6c0fb

	// Address encoding magics
	localParams.PubKeyHashAddrID = []byte{48} // base58 prefix: L
	localParams.ScriptHashAddrID = []byte{50} // base58 prefix: M
	localParams.Bech32HRPSegwit = "ltc"

	if !chaincfg.IsRegistered(&chaincfg.MainNetParams) {
		chaincfg.RegisterBitcoinParams()
	}
	if !chaincfg.IsRegistered(&localParams) {
		err := chaincfg.Register(&localParams)
		if err != nil {
			panic(err)
		}
	}
	return &localParams
}

// GetBTGMainNetParams BTG
func GetBTGMainNetParams() *chaincfg.Params {
	localParams := chaincfg.MainNetParams
	localParams.Net = 0x446d47e1

	// Address encoding magics
	localParams.PubKeyHashAddrID = []byte{38} // base58 prefix: G
	localParams.ScriptHashAddrID = []byte{23} // base58 prefix: A

	// Human-readable part for Bech32 encoded segwit addresses, as defined in
	// BIP 173.
	// see https://github.com/satoshilabs/slips/blob/master/slip-0173.md
	localParams.Bech32HRPSegwit = "btg"

	if !chaincfg.IsRegistered(&localParams) {
		err := chaincfg.Register(&localParams)
		if err != nil {
			panic(err)
		}
	}
	return &localParams
}

// GetDASHMainNetParams DASH
func GetDASHMainNetParams() *chaincfg.Params {
	localParams := chaincfg.MainNetParams
	localParams.Net = 0xbd6b0cbf

	// Address encoding magics
	localParams.PubKeyHashAddrID = []byte{76} // base58 prefix: X
	localParams.ScriptHashAddrID = []byte{16} // base58 prefix: 7

	if !chaincfg.IsRegistered(&localParams) {
		err := chaincfg.Register(&localParams)
		if err != nil {
			panic(err)
		}
	}
	return &localParams
}

// GetDGBMainNetParams DGB
func GetDGBMainNetParams() *chaincfg.Params {
	localParams := chaincfg.MainNetParams
	localParams.Net = 0xdab6c3fa

	// Address encoding magics
	localParams.PubKeyHashAddrID = []byte{30} // base58 prefix: D
	localParams.ScriptHashAddrID = []byte{63} // base58 prefix: 3
	localParams.Bech32HRPSegwit = "dgb"

	if !chaincfg.IsRegistered(&localParams) {
		err := chaincfg.Register(&localParams)
		if err != nil {
			panic(err)
		}
	}
	return &localParams
}

// GetDOGEMainNetParams DOGE
func GetDOGEMainNetParams() *chaincfg.Params {
	localParams := chaincfg.MainNetParams
	localParams.Net = 0xc0c0c0c0

	// Address encoding magics
	localParams.PubKeyHashAddrID = []byte{30} // base58 prefix: D
	localParams.ScriptHashAddrID = []byte{22} // base58 prefix: 9

	if !chaincfg.IsRegistered(&localParams) {
		err := chaincfg.Register(&localParams)
		if err != nil {
			panic(err)
		}
	}
	return &localParams
}

// GetQTUMMainNetParams QTUM
func GetQTUMMainNetParams() *chaincfg.Params {
	localParams := chaincfg.MainNetParams
	localParams.Net = 0xf1cfa6d3

	// Address encoding magics
	localParams.PubKeyHashAddrID = []byte{58} // base58 prefix: Q
	localParams.ScriptHashAddrID = []byte{50} // base58 prefix: P
	localParams.Bech32HRPSegwit = "qc"

	if !chaincfg.IsRegistered(&localParams) {
		err := chaincfg.Register(&localParams)
		if err != nil {
			panic(err)
		}
	}
	return &localParams
}

// GetRVNMainNetParams RVN
func GetRVNMainNetParams() *chaincfg.Params {
	localParams := chaincfg.MainNetParams
	localParams.Net = 0x4e564152

	// Address encoding magics
	localParams.PubKeyHashAddrID = []byte{60}  // base58 prefix: R
	localParams.ScriptHashAddrID = []byte{122} // base58 prefix: r

	if !chaincfg.IsRegistered(&localParams) {
		err := chaincfg.Register(&localParams)
		if err != nil {
			panic(err)
		}
	}
	return &localParams
}

// GetZECMainNetParams ZEC
func GetZECMainNetParams() *chaincfg.Params {
	localParams := chaincfg.MainNetParams
	localParams.Net = 0x6427e924

	// Address encoding magics
	localParams.AddressMagicLen = 2
	localParams.PubKeyHashAddrID = []byte{0x1C, 0xB8} // base58 prefix: t1
	localParams.ScriptHashAddrID = []byte{0x1C, 0xBD} // base58 prefix: t3

	if !chaincfg.IsRegistered(&localParams) {
		err := chaincfg.Register(&localParams)
		if err != nil {
			panic(err)
		}
	}
	return &localParams
}

func GuessUtxoCoinAddressType(address string) string {
	match1, _ := regexp.MatchString("^[1-9A-Za-z]{26,35}$", address)
	if match1 {
		if address[0:1] == "1" || address[0:1] == "L" || address[0:1] == "X" || address[0:1] == "G" || address[0:1] == "t1" || address[0:1] == "D" || address[0:1] == "Q" || address[0:1] == "R" {
			return "P2PKH"
		}
		if address[0:1] == "3" || address[0:1] == "M" || address[0:1] == "7" || address[0:1] == "A" || address[0:1] == "t3" || address[0:1] == "9" || address[0:1] == "P" || address[0:1] == "r" {
			return "P2SH"
		}
	}
	cashAddr, _ := regexp.MatchString("^bitcoincash:[0-9a-zA-Z]{42}$", address)
	if cashAddr {
		if address[12:13] == "q" {
			return "P2PKH"
		} else if address[12:13] == "p" {
			return "P2SH"
		} else {
			return ""
		}
	}
	cashAddrMatch1, _ := regexp.MatchString("^q[0-9a-zA-Z]{30,50}$", strings.ToLower(address))
	if cashAddrMatch1 {
		return "P2PKH"
	}
	cashAddrMatch2, _ := regexp.MatchString("^p[0-9a-zA-Z]{30,50}$", strings.ToLower(address))
	if cashAddrMatch2 {
		return "P2SH"
	}

	ecashAddr, _ := regexp.MatchString("^ecash:[0-9a-zA-Z]{42}$", address)
	if ecashAddr {
		if address[6:7] == "q" {
			return "P2PKH"
		} else if address[6:7] == "p" {
			return "P2SH"
		} else {
			return ""
		}
	}

	if len(address) == 40 {
		return "P2WPKH"
	}

	if len(address) == 64 {
		return "P2WSH"
	}

	match2, _ := regexp.MatchString("^bc1[0-9a-zA-Z]{11,71}$", strings.ToLower(address))
	if match2 {
		return "P2WSH"
	}

	match3, _ := regexp.MatchString("^btg1[0-9a-zA-Z]{11,71}$", strings.ToLower(address))
	if match3 {
		return "P2WSH"
	}

	match4, _ := regexp.MatchString("^ltc1[0-9a-zA-Z]{11,71}$", strings.ToLower(address))
	if match4 {
		return "P2WSH"
	}

	match5, _ := regexp.MatchString("^dgb1[0-9a-zA-Z]{11,71}$", strings.ToLower(address))
	if match5 {
		return "P2WSH"
	}

	match6, _ := regexp.MatchString("^qc1[0-9a-zA-Z]{11,71}$", strings.ToLower(address))
	if match6 {
		return "P2WSH"
	}

	return ""
}

func ConvertCashAddressToLegacy(cashAddr string) (legacy string, err error) {
	addr, err := bchutil.DecodeAddress(cashAddr, GetBTCMainNetParams())
	if err != nil {
		return "", err
	}
	addrType := GuessUtxoCoinAddressType(cashAddr)
	switch addrType {
	case "P2PKH":
		legacyAddr, err := btcutil.NewAddressPubKey(addr.ScriptAddress(), GetBTCMainNetParams())
		if err != nil {
			return "", err
		}
		legacy = legacyAddr.EncodeAddress()

	case "P2SH":
		legacyAddr, err := btcutil.NewAddressScriptHashFromHash(addr.ScriptAddress(), GetBTCMainNetParams())
		if err != nil {
			return "", err
		}
		legacy = legacyAddr.EncodeAddress()
	}

	return legacy, nil
}

func IsCashAddress(addr string) bool {
	if strings.HasPrefix(addr, "bitcoincash:") || strings.HasPrefix(addr, "ecash:") {
		return true
	}

	if strings.HasPrefix(addr, "q") || strings.HasPrefix(addr, "p") {
		return true
	}

	return false
}
