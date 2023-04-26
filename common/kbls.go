package common

import (
	kbls "github.com/kilic/bls12-381"
)

// IETF signature draft v4:
// https://datatracker.ietf.org/doc/html/draft-irtf-cfrg-bls-signature-04
//

var domain = []byte("BLS_SIG_BLS12381G2_XMD:SHA-256_SSWU_RO_POP_")

type Pubkey kbls.PointG1

func (pub *Pubkey) Serialize() (out [48]byte) {
	copy(out[:], kbls.NewG1().ToCompressed((*kbls.PointG1)(pub)))
	return
}

func (pub *Pubkey) Deserialize(in *[48]byte) error {
	p, err := kbls.NewG1().FromCompressed(in[:])
	if err != nil {
		return err
	}
	*pub = (Pubkey)(*p)
	return nil
}

type Signature kbls.PointG2

func (sig *Signature) Serialize() (out [96]byte) {
	copy(out[:], kbls.NewG2().ToCompressed((*kbls.PointG2)(sig)))
	return
}

func (sig *Signature) Deserialize(in *[96]byte) error {
	// includes sub-group check
	p, err := kbls.NewG2().FromCompressed(in[:])
	if err != nil {
		return err
	}
	*sig = (Signature)(*p)
	return nil
}

// The coreVerify algorithm checks that a signature is valid for the octet string message under the public key PK.
func coreVerify(pk *Pubkey, message []byte, signature *Signature) bool {
	// 1. R = signature_to_point(signature)
	R := (*kbls.PointG2)(signature)
	// 2. If R is INVALID, return INVALID
	// 3. If signature_subgroup_check(R) is INVALID, return INVALID
	// 4. If KeyValidate(PK) is INVALID, return INVALID
	// steps 2-4 are part of bytes -> *Signature deserialization
	if (*kbls.G2)(nil).IsZero(R) {
		// KeyValidate is assumed through deserialization of Pubkey and Signature,
		// but the identity pubkey/signature case is not part of that, thus verify here.
		return false
	}

	// 5. xP = pubkey_to_point(PK)
	xP := (*kbls.PointG1)(pk)
	// 6. Q = hash_to_point(message)
	Q, err := kbls.NewG2().HashToCurve(message, domain)
	if err != nil {
		// e.g. when the domain is too long. Maybe change to panic if never due to a usage error?
		return false
	}
	// 7. C1 = pairing(Q, xP)
	eng := kbls.NewEngine()
	eng.AddPair(xP, Q)
	// 8. C2 = pairing(R, P)
	P := &kbls.G1One
	eng.AddPairInv(P, R) // inverse, optimization to mul with inverse and check equality to 1
	// 9. If C1 == C2, return VALID, else return INVALID
	return eng.Check()
}

// The Verify algorithm checks an aggregated signature over several (PK, message) pairs.
func Verify(pk *Pubkey, message []byte, signature *Signature) bool {
	return coreVerify(pk, message, signature)
}
