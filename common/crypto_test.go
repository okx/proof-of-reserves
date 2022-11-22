package common

import (
	"crypto/sha256"
	"github.com/btcsuite/btcd/btcutil"
	"github.com/btcsuite/btcd/chaincfg"
	"github.com/btcsuite/btcd/txscript"
	"testing"
)

func TestETHVerifySignature(t *testing.T) {
	addr := "0x4ce08ffc090f5c54013c62efe30d62e6578e738d"
	msg := "I am an OKX address"
	sign := "0x1a3ca543e08a1402baafd91aa56c478338322a2fa08f2fa2edc36b44a94a49806aab8ab444be066123eac18e9b206b820eca948bc7967fafcbfea1c87efb1a8c1c"
	if err := VerifyETH(addr, msg, sign); err != nil {
		t.Errorf(err.Error())
	}
}
func TestSOLVerifySignature(t *testing.T) {
	pub := "7bzoTJhZmpU1vQVjN63fQ3iVYmWCVgQh1sYSqsjuapU9"
	msg := "hello world"
	sign := "oaxtr4HVifyZwYfCbqrfheo7vTpNXQu8kuwWpHJ57BaTiBGkiysaA8YjnkUTc2fJWHPNbndcwRtBXEZ9V8gKsnq"
	if err := VerifySol(pub, msg, sign); err != nil {
		t.Errorf(err.Error())
	}
}

func TestTRXVerifySignature(t *testing.T) {
	addr := "TEjxQjU3CxkFrSDcPfHwZXSuPpCpdQ27NJ"
	msg := "hello world"
	sign := "0xcd1e3903dc047ea881f7da1647fa3372f37ee6a1cf0726477a20e267408af43f3f9c3a43f7f15e6bf674c9f0776866b6d6a770ce998b29cc03f11f2cb98df5821c"
	if err := VerifyTRX(addr, msg, sign); err != nil {
		t.Errorf(err.Error())
	}
}

func TestBTCVerifySignature(t *testing.T) {
	addr := "18WpcobD4TPJSfyeFUjJrtavYHRxiEs7gD"
	msg := "I am an OKX address"
	sign := "H5iQGmrUlmVGhxjq/Yu88najhwKS8ZiRv4YHe0n4L+Q7e+S10TsAjq7mIYxLnpDq/078MkyOsOw7luhGVBH24Hc="
	if err := VerifyBTC(addr, msg, sign); err != nil {
		t.Errorf("invalid address")
	}
}

func TestBTCWitnessVerifySignature(t *testing.T) {
	msg := "I am an OKX address"
	addr := "bc1q6pue5m5a0vf27cdpns23r5pku6evzh93jpes4mrypmy4ya8xyutqug6wn5"
	script := MustDecode("522102f875caa8d0916852f2fe4edbdd42d71ccd223b2dca9710b9967236877d0d7d882103a7c1c409c521b41156e0d0af386e53ca7c525f905122d15aabc9b44cd042662b210384fe09c13065e47582b8dd915c7a05e39c5dc3e0d241b2ee16e964052c96273e53ae")
	h := sha256.New()
	h.Write(script)
	witnessProg := h.Sum(nil)
	addressWitnessScriptHash, err := btcutil.NewAddressWitnessScriptHash(witnessProg, &chaincfg.MainNetParams)
	if err != nil {
		t.Errorf("invalid address")
		t.Fail()
	}
	if addressWitnessScriptHash.EncodeAddress() != addr {
		t.Errorf("invalid address")
		t.Fail()
	}
	addr1, err := SigToAddrBTC(msg, "IEpZHHUHLP7y3Uz+XoQgGug6H4BygnLco/57/tnI5FpUNmSJoOhx4mIQcEFtD1VVjKPK9hP57gE7jI/iV7pj7vc=")
	if err != nil {
		t.Errorf("invalid address")
		t.Fail()
	}
	//addr2, err := SigToAddrBTC(msg, "H5DSCKyIOGnCUk7uust3fH+QN/QdSQf1FC3ddZf/nY3xZk5bEKrw5Siy/zF7gWxoZn8/BVA635wbSgUH4JZkUFo=")
	//if err != nil {
	//	t.Errorf("invalid address")
	//	t.Fail()
	//}
	m := map[string]struct{}{addr1: {}}
	_, addrs, _, err := txscript.ExtractPkScriptAddrs(script, &chaincfg.MainNetParams)
	if err != nil {
		t.Errorf("invalid script")
		t.Fail()
	}
	for _, v := range addrs {
		delete(m, v.EncodeAddress())
	}
	if len(m) > 1 {
		t.Errorf("invalid script and address")
		t.Fail()
	}
}
