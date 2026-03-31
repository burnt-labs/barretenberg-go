//go:build cgo
// +build cgo

package barretenberg

import (
	"testing"
)

// TestProofBytesImmutable verifies that mutating the slice returned by Proof.Bytes()
// does not alter the Proof's internal state or change what Proof.Hex() returns.
func TestProofBytesImmutable(t *testing.T) {
	data := make([]byte, 1024)
	for i := range data {
		data[i] = byte(i % 256)
	}

	proof, err := ParseProof(data)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	hexBefore := proof.Hex()

	// Obtain bytes and corrupt them entirely
	b := proof.Bytes()
	for i := range b {
		b[i] = 0xFF
	}

	hexAfter := proof.Hex()
	if hexBefore != hexAfter {
		t.Error("Proof.Hex() changed after mutating Proof.Bytes() result — internal state was corrupted")
	}

	// Verify Size() is also unaffected
	if proof.Size() != len(data) {
		t.Errorf("Proof.Size() = %d after mutation, expected %d", proof.Size(), len(data))
	}
}

// TestProofBytesReturnsCopy verifies that successive calls to Proof.Bytes() return
// independent copies (modifying one does not affect subsequent calls).
func TestProofBytesReturnsCopy(t *testing.T) {
	data := make([]byte, 1024)
	for i := range data {
		data[i] = byte(i % 256)
	}

	proof, err := ParseProof(data)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	b1 := proof.Bytes()
	b1[0] = 0xFF // Mutate first copy

	b2 := proof.Bytes()
	if b2[0] == 0xFF {
		t.Error("second Proof.Bytes() call returned slice aliasing first call's result")
	}
	if b2[0] != data[0] {
		t.Errorf("second Proof.Bytes()[0] = 0x%02x, expected original 0x%02x", b2[0], data[0])
	}
}

// TestVerificationKeyBytesImmutable verifies that mutating the slice returned by
// VerificationKey.Bytes() does not alter VerificationKey.Hex().
// Requires real test vectors to parse a valid vkey; skipped if not available.
func TestVerificationKeyBytesImmutable(t *testing.T) {
	vkeyData := loadTestVector(t, vkeyFile)

	vk, err := ParseVerificationKey(vkeyData)
	if err != nil {
		t.Fatalf("failed to parse vkey: %v", err)
	}
	defer vk.Close()

	hexBefore := vk.Hex()

	// Obtain bytes and corrupt them entirely
	b := vk.Bytes()
	for i := range b {
		b[i] = 0xFF
	}

	hexAfter := vk.Hex()
	if hexBefore != hexAfter {
		t.Error("VerificationKey.Hex() changed after mutating VerificationKey.Bytes() result — internal state was corrupted")
	}
}

// TestVerificationKeyBytesReturnsCopy verifies successive Bytes() calls are independent.
func TestVerificationKeyBytesReturnsCopy(t *testing.T) {
	vkeyData := loadTestVector(t, vkeyFile)

	vk, err := ParseVerificationKey(vkeyData)
	if err != nil {
		t.Fatalf("failed to parse vkey: %v", err)
	}
	defer vk.Close()

	b1 := vk.Bytes()
	original0 := b1[0]
	b1[0] = original0 ^ 0xFF // Flip bits

	b2 := vk.Bytes()
	if b2[0] != original0 {
		t.Errorf("second VerificationKey.Bytes()[0] = 0x%02x, expected original 0x%02x — slice aliased first call", b2[0], original0)
	}
}
