package main

import "testing"

// Real, public ETH address + signature reused from example/okx_por_example.csv (this address
// self-signs the OKX message). No private keys or generated crypto — just a fixed public vector.
const (
	okxMsg  = "I am an OKX address"
	ethAddr = "0x0cdcdb19a857c2ac24818ca4fdfe38cce071483e"
	ethSig  = "0x07f19879aa28d51c97cddfdfecffe7ed96525545d041aee4f4386b0bf4c1a26924b637fb02ccbb97305c13daa51a0f50b8896fb25ecbaf60020cde920d227a221b"
)

// FR-6 / AC-1: the same ETH row verifies under both the legacy 11-column layout (off=0) and the
// new 12-column layout with a Type column inserted after coin (off=1), proving the header-adaptive
// column mapping selects the right fields in both formats.
func TestHandleParsesBothFormats(t *testing.T) {
	oldRow := "ETH,ETH,20914735," + ethAddr + ",0.0589," + okxMsg + "," + ethSig + ",,"
	if _, ok := handle(0, oldRow, 0); !ok {
		t.Fatalf("legacy 11-column row should verify")
	}
	newRow := "ETH,Non Staking,ETH,20914735," + ethAddr + ",0.0589," + okxMsg + "," + ethSig + ",,,,"
	if _, ok := handle(1, newRow, 1); !ok {
		t.Fatalf("12-column row should verify (column shift handled)")
	}
}

// A staking-typed row (address = display-only validator pubkey, EOA1 populated) verifies through
// the EXISTING EVM eoa1 owner-mode branch — no staking-specific verification code is added.
func TestHandleStakingRowUsesExistingEoaBranch(t *testing.T) {
	validatorPub := "0x" + "ab" // display-only placeholder; never used for verification
	nativeRow := "ETH,Native ETH Staking,ETH,20914735," + validatorPub + ",32," + okxMsg + "," + ethSig + ",,," + ethAddr + ","
	if _, ok := handle(2, nativeRow, 1); !ok {
		t.Fatalf("Native ETH Staking row should verify via the existing eoa1 branch")
	}
}

// A row with fewer columns than the detected format fails cleanly (no panic).
func TestHandleTooFewColumns(t *testing.T) {
	if _, ok := handle(3, "ETH,Non Staking,ETH", 1); ok {
		t.Fatalf("a short row must fail, not panic")
	}
}

// A Pearl 2-of-2 row (signatures and tapscript from okcoin-signature PearlPorTest) routes by the
// Network column to common.VerifyPearl; swapping the two signatures must fail.
func TestHandlePearlRow(t *testing.T) {
	const (
		addr   = "prl1pg8uak7ghkg84y3fdzqrjf4s0e8md6l5galn4j8q40pc98k3nf62qaay2u2"
		script = "208ce13849b588985e19f2fffe34c2c98150efb70285b17557da647379b188a1cfac200203711eaf703c17670691044de6cfc706b5fe091e78b80f6b6bb4b578151996ba529c"
		sig1   = "2959c03e595c6c86ce1cf425c3a285e46444861e4afa494f30832b5064e3ae5ca73e0238d20cfd8df627f974b2fb23b81a2b0f613b2da85d0f96acd3e3212b74"
		sig2   = "4adfe80f213012e8ed700802726b8a96e1072d35973943372039bc510a832b1bc00fbebcecf461305748095b75f9249b908bdf2217f88da3c5311b83b7203c39"
	)
	row := "PRL,Non Staking,PEARL,1," + addr + ",1.5,hello world," + sig1 + "," + sig2 + "," + script + ",,"
	if _, ok := handle(4, row, 1); !ok {
		t.Fatalf("Pearl row should verify")
	}
	swapped := "PRL,Non Staking,PEARL,1," + addr + ",1.5,hello world," + sig2 + "," + sig1 + "," + script + ",,"
	if _, ok := handle(5, swapped, 1); ok {
		t.Fatalf("Pearl row with swapped signatures must fail")
	}
}
