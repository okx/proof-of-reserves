package common

import (
	"crypto/ed25519"
	"crypto/rand"
	"encoding/hex"
	"strings"
	"testing"
)

// These tests cover the APTOS key-rotation "owner mode" at the crypto layer: verifying an
// ed25519 signature and comparing sha3_256(pubkey||0x00) against eoa1 (the current
// authentication key) instead of the claimed address. Owner mode is implemented in
// cmd/verifyaddress by passing eoa1 in the address slot of VerifyEd25519Coin, so these
// tests exercise that exact call shape. All key material is generated locally (no real
// production data).

const aptosOwnerMsg = "I am an OKX address"

// aptosAuthKey returns "0x" + hex(sha3_256(pub||0x00)) — the Aptos authentication key.
func aptosAuthKey(pub ed25519.PublicKey) string {
	buf := append(append([]byte{}, pub...), 0x0)
	return "0x" + hex.EncodeToString(Sha256Hash(buf))
}

func signOKXEd25519(t *testing.T, priv ed25519.PrivateKey, msg string) string {
	t.Helper()
	return "0x" + hex.EncodeToString(ed25519.Sign(priv, HashEd25519Msg(OKXMessageSignatureHeader, msg)))
}

func flipLastHexNibble(s string) string {
	if len(s) == 0 {
		return s
	}
	repl := byte('0')
	if s[len(s)-1] == '0' {
		repl = '1'
	}
	return s[:len(s)-1] + string(repl)
}

// AC-1.2 / AC-1.3 / AC-1.4
func TestAptosOwnerModeVerifyEd25519Coin(t *testing.T) {
	pub, priv, err := ed25519.GenerateKey(rand.Reader)
	if err != nil {
		t.Fatalf("GenerateKey: %v", err)
	}
	pubHex := "0x" + hex.EncodeToString(pub)
	authKey := aptosAuthKey(pub)
	sign := signOKXEd25519(t, priv, aptosOwnerMsg)

	// AC-1.2: valid sig + sha3_256(pub||0x00) == eoa1 -> pass (claimed address irrelevant).
	if err := VerifyEd25519Coin("APTOS", authKey, aptosOwnerMsg, sign, pubHex); err != nil {
		t.Fatalf("AC-1.2 owner-mode should pass, got: %v", err)
	}

	// AC-1.3: valid sig but eoa1 != derived auth_key -> "recovery address not match".
	wrongOwner := "0xdeadbeef00000000000000000000000000000000000000000000000000000000"
	if err := VerifyEd25519Coin("APTOS", wrongOwner, aptosOwnerMsg, sign, pubHex); err == nil {
		t.Fatalf("AC-1.3 should fail when eoa1 != derived auth_key")
	} else if !strings.Contains(err.Error(), "recovery address not match") {
		t.Fatalf("AC-1.3 wrong error: %v", err)
	}

	// AC-1.4: eoa1 matches but signature tampered -> "ED25519 signature verification failed".
	if err := VerifyEd25519Coin("APTOS", authKey, aptosOwnerMsg, flipLastHexNibble(sign), pubHex); err == nil {
		t.Fatalf("AC-1.4 should fail on tampered signature")
	} else if !strings.Contains(err.Error(), "ED25519 signature verification failed") {
		t.Fatalf("AC-1.4 wrong error: %v", err)
	}
}

// AC-3.1: a leading-zero auth_key matches in BOTH the full 64-hex form and the
// leading-zero-stripped form (reusing the existing double-form compare; no new code).
func TestAptosOwnerModeLeadingZero(t *testing.T) {
	var pub ed25519.PublicKey
	var priv ed25519.PrivateKey
	var authKeyFull string
	found := false
	for i := 0; i < 200000; i++ {
		p, s, err := ed25519.GenerateKey(rand.Reader)
		if err != nil {
			t.Fatalf("GenerateKey: %v", err)
		}
		raw := Sha256Hash(append(append([]byte{}, p...), 0x0))
		if raw[0] == 0x00 {
			pub, priv, authKeyFull = p, s, "0x"+hex.EncodeToString(raw)
			found = true
			break
		}
	}
	if !found {
		t.Skip("no leading-zero auth_key found (improbable)")
	}
	pubHex := "0x" + hex.EncodeToString(pub)
	sign := signOKXEd25519(t, priv, aptosOwnerMsg)

	// full form (with the leading zero present)
	if err := VerifyEd25519Coin("APTOS", authKeyFull, aptosOwnerMsg, sign, pubHex); err != nil {
		t.Fatalf("AC-3.1 full-form leading-zero auth_key should pass, got: %v", err)
	}
	// stripped form (leading zeros after 0x removed)
	stripped := "0x" + strings.TrimLeft(authKeyFull[2:], "0")
	if err := VerifyEd25519Coin("APTOS", stripped, aptosOwnerMsg, sign, pubHex); err != nil {
		t.Fatalf("AC-3.1 stripped-form auth_key should pass, got: %v (stripped=%s)", err, stripped)
	}
}
