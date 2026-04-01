//go:build cgo
// +build cgo

package barretenberg

import (
	"errors"
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
