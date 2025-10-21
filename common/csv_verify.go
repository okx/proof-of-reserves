package common

import (
	"encoding/json"
	"fmt"
	"strings"
	"testing"
)

// Verify single line data (single-threaded version, only handle StarkNet)
func verifyCSVLineStarknetOnly(coin, addr, msg, sign1, sign2, publicKey, digitalAsset string, lineNumber int, t *testing.T) (bool, string, string) {
	// Check if it's StarkNet coin, skip if not
	coinType, exists := PorCoinTypeMap[coin]
	if !exists || coinType != StarkCoinType {
		return true, coin, ""
	}

	return verifyCSVLineInternal(coin, addr, msg, sign1, sign2, publicKey, "", "", digitalAsset, lineNumber, t)
}

// Verify single line data (multithreaded version, skip StarkNet)
func verifyCSVLineMultithread(coin, addr, msg, sign1, sign2, publicKey, owner1, owner2, digitalAsset string, lineNumber int, t *testing.T) (bool, string, string) {
	// Check if it's StarkNet coin, skip multithreaded verification if yes
	coinType, exists := PorCoinTypeMap[coin]
	if exists && coinType == StarkCoinType {
		return true, coin, ""
	}

	return verifyCSVLineInternal(coin, addr, msg, sign1, sign2, publicKey, owner1, owner2, digitalAsset, lineNumber, t)
}

// Internal verification logic (shared)
func verifyCSVLineInternal(coin, addr, msg, sign1, sign2, publicKey, owner1, owner2, digitalAsset string, lineNumber int, t *testing.T) (bool, string, string) {

	// Check required fields
	if addr == "" || msg == "" || sign1 == "" {
		errorMsg := fmt.Sprintf("Missing required parameters (digitalAsset:%s, network:%s, addr:%s)", digitalAsset, coin, addr)
		t.Logf("Line %d: %s", lineNumber, errorMsg)
		return false, coin, errorMsg
	}

	// Get coin verification type
	coinType, exists := PorCoinTypeMap[coin]
	if !exists {
		// If coin is not in mapping table, try using default ECDSA verification
		coinType = EcdsaCoinType
	}

	// Verify according to coin type
	var err error
	switch coinType {
	case EvmCoinTye:
		// For EVM coins, if there are owner fields, use owner for signature verification
		if owner1 != "" && owner1 != "null" {
			// Use signature1 + owner1 for verification
			err = VerifyEvmCoin(coin, owner1, msg, sign1)

			// If there's a second owner, both must pass verification
			if owner2 != "" && owner2 != "null" && sign2 != "" && sign2 != "null" {
				err2 := VerifyEvmCoin(coin, owner2, msg, sign2)

				// Record verification results
				if err == nil && err2 == nil {
					t.Logf("Line %d: Dual owner verification both successful (digitalAsset:%s, network:%s, owner1:%s, owner2:%s)", lineNumber, digitalAsset, coin, owner1, owner2)
				} else if err != nil && err2 != nil {
					// Both failed
					t.Logf("Line %d: EVM contract address dual verification both failed (digitalAsset:%s, network:%s, addr:%s, owner1:%s, owner2:%s)", lineNumber, digitalAsset, coin, addr, owner1, owner2)
					err = fmt.Errorf("owner1 verification failed: %v, owner2 verification failed: %v", err, err2)
				} else if err != nil {
					// owner1 failed, owner2 succeeded
					t.Logf("Line %d: EVM contract address owner1 verification failed, owner2 succeeded (digitalAsset:%s, network:%s, addr:%s, owner1:%s, owner2:%s)", lineNumber, digitalAsset, coin, addr, owner1, owner2)
					err = fmt.Errorf("owner1 verification failed: %v", err)
				} else {
					// owner1 succeeded, owner2 failed
					t.Logf("Line %d: EVM contract address owner1 verification succeeded, owner2 failed (digitalAsset:%s, network:%s, addr:%s, owner1:%s, owner2:%s)", lineNumber, digitalAsset, coin, addr, owner1, owner2)
					err = fmt.Errorf("owner2 verification failed: %v", err2)
				}
			} else {
				// Single owner case
				if err != nil {
					t.Logf("Line %d: EVM contract address owner1 verification failed (digitalAsset:%s, network:%s, addr:%s, owner1:%s)", lineNumber, digitalAsset, coin, addr, owner1)
				}
			}
		} else {
			// No owner field, use original address for verification
			err = VerifyEvmCoin(coin, addr, msg, sign1)
		}
	case EcdsaCoinType:
		// For XRP and other coins using ECDSA, use public key verification if available
		if publicKey != "" && publicKey != "null" {
			err = VerifyEcdsaCoinWithPub(msg, sign1, publicKey)
		} else {
			err = VerifyEcdsaCoin(coin, addr, msg, sign1)
		}
	case Ed25519CoinType:
		if publicKey == "" || publicKey == "null" {
			errorMsg := fmt.Sprintf("ED25519 coin %s missing public key (digitalAsset:%s)", coin, digitalAsset)
			t.Logf("Line %d: %s", lineNumber, errorMsg)
			return false, coin, errorMsg
		}
		err = VerifyEd25519Coin(coin, addr, msg, sign1, publicKey)
	case TrxCoinType:
		err = VerifyTRX(addr, msg, sign1)
	case BethCoinType:
		err = VerifyBETH(addr, msg, sign1)
	case UTXOCoinType:
		err = VerifyUtxoCoin(coin, addr, msg, sign1, sign2, publicKey)
	case StarkCoinType:
		if publicKey == "" || publicKey == "null" {
			errorMsg := fmt.Sprintf("STARK coin %s missing public key (digitalAsset:%s)", coin, digitalAsset)
			t.Logf("Line %d: %s", lineNumber, errorMsg)
			return false, coin, errorMsg
		}
		err = VerifyStarkCoin(coin, addr, msg, sign1, publicKey)
	case EOSCoinType:
		if publicKey == "" || publicKey == "null" {
			errorMsg := fmt.Sprintf("EOS coin %s missing public key (digitalAsset:%s)", coin, digitalAsset)
			t.Logf("Line %d: %s", lineNumber, errorMsg)
			return false, coin, errorMsg
		}

		cleanKey := strings.TrimSpace(publicKey)
		if strings.HasPrefix(cleanKey, "\"") && strings.HasSuffix(cleanKey, "\"") {
			cleanKey = cleanKey[1 : len(cleanKey)-1]
			cleanKey = strings.ReplaceAll(cleanKey, "\"\"", "\"") // 处理CSV转义的双引号
		}

		if strings.HasPrefix(cleanKey, "{") {
			var pubKeys map[string]string
			if json.Unmarshal([]byte(cleanKey), &pubKeys) == nil {
				err1 := VerifyEOSCoin(coin, addr, msg, sign1, pubKeys["publicKey1"])
				err2 := VerifyEOSCoin(coin, addr, msg, sign2, pubKeys["publicKey2"])
				if err1 != nil || err2 != nil {
					err = fmt.Errorf("EOS dual signature failed: sig1=%v, sig2=%v", err1, err2)
				}
			} else {
				err = fmt.Errorf("invalid JSON public key format")
			}
		} else {
			err = VerifyEOSCoin(coin, addr, msg, sign1, publicKey)
		}
	default:
		errorMsg := fmt.Sprintf("Unsupported coin type %s (digitalAsset:%s, network:%s)", coinType, digitalAsset, coin)
		t.Logf("Line %d: %s", lineNumber, errorMsg)
		return false, coin, errorMsg
	}

	if err != nil {
		errorMsg := fmt.Sprintf("Verification failed: %v", err)
		t.Logf("Line %d verification failed: %s (digitalAsset:%s, network:%s, addr:%s, error:%v)", lineNumber, coin, digitalAsset, coin, addr, err)
		return false, coin, errorMsg
	}

	return true, coin, ""
}
