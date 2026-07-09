package main

import "testing"

// Real, public ETH address + signature reused from example/okx_por_example.csv (this address
// self-signs the OKX message). No private keys or generated crypto — just a fixed public vector.
const (
	okxMsg  = "I am an OKX address"
	ethAddr = "0x0cdcdb19a857c2ac24818ca4fdfe38cce071483e"
	ethSig  = "0x07f19879aa28d51c97cddfdfecffe7ed96525545d041aee4f4386b0bf4c1a26924b637fb02ccbb97305c13daa51a0f50b8896fb25ecbaf60020cde920d227a221b"
)

// FR-6: format is detected from the header — a "Type" column marks the 12-column layout.
func TestDetectFormatOffset(t *testing.T) {
	oldHeader := []string{"coin", "Network", "Snapshot Height", "address", "amount", "message", "signature1", "signature2", "redeem script/ public key", " eoa1", " eoa2"}
	if got := detectFormatOffset(oldHeader); got != 0 {
		t.Errorf("legacy header offset = %d, want 0", got)
	}
	newHeader := []string{"coin", "Type", "Network", "Snapshot Height", "address", "amount", "message", "signature1", "signature2", "redeem script/public key", "EOA1", "EOA2"}
	if got := detectFormatOffset(newHeader); got != 1 {
		t.Errorf("12-column header offset = %d, want 1", got)
	}
}

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
