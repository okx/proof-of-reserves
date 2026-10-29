package common

import (
	"errors"
	"strings"
	"testing"
)

// PoR results produced by the OKX signer (okcoin-signature PearlPorTest, hot
// and cold 2-of-2 on Pearl mainnet). Each signatureN also verifies on its own
// against tapscript keyN over the go-parent-sdk coin/pearl MessageHash.
const (
	pearlMsg = "hello world"

	pearlHotAddr   = "prl1pg8uak7ghkg84y3fdzqrjf4s0e8md6l5galn4j8q40pc98k3nf62qaay2u2"
	pearlHotScript = "208ce13849b588985e19f2fffe34c2c98150efb70285b17557da647379b188a1cfac200203711eaf703c17670691044de6cfc706b5fe091e78b80f6b6bb4b578151996ba529c"
	pearlHotSig1   = "2959c03e595c6c86ce1cf425c3a285e46444861e4afa494f30832b5064e3ae5ca73e0238d20cfd8df627f974b2fb23b81a2b0f613b2da85d0f96acd3e3212b74"
	pearlHotSig2   = "4adfe80f213012e8ed700802726b8a96e1072d35973943372039bc510a832b1bc00fbebcecf461305748095b75f9249b908bdf2217f88da3c5311b83b7203c39"

	pearlColdAddr   = "prl1p45f5jsyxml3p89dv9n0kx7wdvlcux87ehn90k3x9a9xyw5tpn6vsk9lg6a"
	pearlColdScript = "20fc94159c81d848a896d90404f8815aa7ffd61748fdc86c67b294ce65864593e7ac20da326a29ec44da294786c9507bff5bb467800be3e9ec4814e625e5995cb7d765ba529c"
	pearlColdSig1   = "056d0c63b39638bb7be89ce907df1fe156485682d075272c317a724e459f86e8cfa46fc3205674ca79b294b37b5493678167b4fc54aa0b5c0784b99881cd9147"
	pearlColdSig2   = "da6e630ef13e05e1e8ebf1b359a6e3bdd0b7e25241b6b61d22e1c1945c0a038377c9e4148ef0d3dd08227063ec85e9267720713e3cacbf6b9ba400aa51ce12ef"
)

func TestVerifyPearl(t *testing.T) {
	if err := VerifyPearl(pearlHotAddr, pearlMsg, pearlHotSig1, pearlHotSig2, pearlHotScript); err != nil {
		t.Fatalf("hot proof should verify: %v", err)
	}
	if err := VerifyPearl(pearlColdAddr, pearlMsg, pearlColdSig1, pearlColdSig2, pearlColdScript); err != nil {
		t.Fatalf("cold proof should verify: %v", err)
	}
	// Bech32 is case-insensitive; reports may upper-case addresses.
	if err := VerifyPearl(strings.ToUpper(pearlHotAddr), pearlMsg, pearlHotSig1, pearlHotSig2, "0x"+pearlHotScript); err != nil {
		t.Fatalf("upper-case address and 0x script should verify: %v", err)
	}
}

func TestVerifyPearl_Rejects(t *testing.T) {
	cases := []struct {
		name                          string
		addr, msg, sig1, sig2, script string
		want                          error
	}{
		{"wrong message", pearlHotAddr, "I am NOT an OKX address", pearlHotSig1, pearlHotSig2, pearlHotScript, ErrPearlSignature},
		{"swapped signatures", pearlHotAddr, pearlMsg, pearlHotSig2, pearlHotSig1, pearlHotScript, ErrPearlSignature},
		{"other address's signatures", pearlHotAddr, pearlMsg, pearlColdSig1, pearlColdSig2, pearlHotScript, ErrPearlSignature},
		{"missing signature2", pearlHotAddr, pearlMsg, pearlHotSig1, "", pearlHotScript, ErrPearlSignature},
		{"null signature2", pearlHotAddr, pearlMsg, pearlHotSig1, "null", pearlHotScript, ErrPearlSignature},
		{"one signature twice", pearlHotAddr, pearlMsg, pearlHotSig1, pearlHotSig1, pearlHotScript, ErrPearlSignature},
		{"short signature", pearlHotAddr, pearlMsg, pearlHotSig1[:126], pearlHotSig2, pearlHotScript, ErrPearlSignature},
		{"message over 80 bytes", pearlHotAddr, strings.Repeat("a", 81), pearlHotSig1, pearlHotSig2, pearlHotScript, ErrPearlSignature},
		{"script of another address", pearlHotAddr, pearlMsg, pearlHotSig1, pearlHotSig2, pearlColdScript, nil},
		{"keys swapped in script", pearlHotAddr, pearlMsg, pearlHotSig1, pearlHotSig2,
			"20" + pearlHotScript[70:134] + "ac20" + pearlHotScript[2:66] + "ba529c", nil},
		{"not a 2-of-2 leaf", pearlHotAddr, pearlMsg, pearlHotSig1, pearlHotSig2, pearlHotScript[:136] + "529d", ErrPearlScript},
		{"duplicate keys", pearlHotAddr, pearlMsg, pearlHotSig1, pearlHotSig2,
			"20" + pearlHotScript[2:66] + "ac20" + pearlHotScript[2:66] + "ba529c", ErrPearlScript},
		{"empty script", pearlHotAddr, pearlMsg, pearlHotSig1, pearlHotSig2, "", ErrPearlScript},
		{"bitcoin address", "bc1pg8uak7ghkg84y3fdzqrjf4s0e8md6l5galn4j8q40pc98k3nf62qxxxxxx", pearlMsg, pearlHotSig1, pearlHotSig2, pearlHotScript, nil},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			err := VerifyPearl(c.addr, c.msg, c.sig1, c.sig2, c.script)
			if err == nil {
				t.Fatal("expected an error")
			}
			if c.want != nil && !errors.Is(err, c.want) {
				t.Fatalf("expected %v, got %v", c.want, err)
			}
		})
	}
}

// Both report flows route Pearl rows through the UTXO scheme by the Network column.
func TestVerifyPearl_Dispatch(t *testing.T) {
	ok, _, msg := verifyCSVLineInternal("PEARL", pearlHotAddr, pearlMsg, pearlHotSig1, pearlHotSig2, pearlHotScript, "", "", "PRL", 1, t)
	if !ok {
		t.Fatalf("CSV flow should verify the Pearl row: %s", msg)
	}
	ok, _, _ = verifyCSVLineInternal("PRL", pearlHotAddr, pearlMsg, pearlHotSig2, pearlHotSig1, pearlHotScript, "", "", "PRL", 1, t)
	if ok {
		t.Fatal("CSV flow accepted swapped signatures")
	}
}
