//go:build cgo
// +build cgo

package barretenberg

import (
	"errors"
	"strings"
	"testing"
)

// TestProofNilReceiver verifies that all Proof methods return zero-values when
// called on a nil receiver rather than panicking.
func TestProofNilReceiver(t *testing.T) {
	var p *Proof

	if got := p.Bytes(); got != nil {
		t.Errorf("nil Proof.Bytes() = %v, expected nil", got)
	}
	if got := p.Hex(); got != "" {
		t.Errorf("nil Proof.Hex() = %q, expected empty string", got)
	}
	if got := p.Size(); got != 0 {
		t.Errorf("nil Proof.Size() = %d, expected 0", got)
	}
}

// TestPublicInputsNilReceiver verifies that all PublicInputs methods return
// zero-values or safe errors when called on a nil receiver rather than panicking.
func TestPublicInputsNilReceiver(t *testing.T) {
	var pi *PublicInputs

	if got := pi.Count(); got != 0 {
		t.Errorf("nil PublicInputs.Count() = %d, expected 0", got)
	}
	if got := pi.Bytes(); got != nil {
		t.Errorf("nil PublicInputs.Bytes() = %v, expected nil", got)
	}
	_, err := pi.Element(0)
	if err == nil {
		t.Error("nil PublicInputs.Element(0): expected error, got nil")
	}
	if !errors.Is(err, ErrInvalidPublicInputs) {
		t.Errorf("nil PublicInputs.Element(0): expected ErrInvalidPublicInputs, got %v", err)
	}
}

// TestVerifyPublicInputCountMismatch verifies that verifying with the wrong
// number of public inputs returns ErrInvalidPublicInputs with a clear message.
// Requires real test vectors; skipped with stub library or missing data.
func TestVerifyPublicInputCountMismatch(t *testing.T) {
	if strings.HasPrefix(Version(), "stub") {
		t.Skip("stub library does not support NumPublicInputs")
	}

	vkeyData := loadTestVector(t, vkeyFile)
	proofData := loadTestVector(t, proofFile)

	vkey, err := ParseVerificationKey(vkeyData)
	if err != nil {
		t.Fatalf("failed to parse vkey: %v", err)
	}
	defer vkey.Close()

	expectedCount, err := vkey.NumPublicInputs()
	if err != nil {
		t.Fatalf("failed to get public input count: %v", err)
	}
	if expectedCount == 0 {
		t.Skip("vkey expects 0 public inputs — cannot test count mismatch with this test vector")
	}

	proof, err := ParseProof(proofData)
	if err != nil {
		t.Fatalf("failed to parse proof: %v", err)
	}

	verifier, err := NewVerifier(vkey)
	if err != nil {
		t.Fatalf("failed to create verifier: %v", err)
	}
	defer verifier.Close()

	// Pass no inputs when vkey expects some — must fail with ErrInvalidPublicInputs
	_, err = verifier.VerifyWithBytes(proof, nil)
	if err == nil {
		t.Fatal("expected error for count mismatch, got nil")
	}
	if !errors.Is(err, ErrInvalidPublicInputs) {
		t.Errorf("expected ErrInvalidPublicInputs for count mismatch, got %v", err)
	}
	t.Logf("correctly rejected: vkey expects %d input(s), got 0 — error: %v", expectedCount, err)
}
