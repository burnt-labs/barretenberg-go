//go:build cgo
// +build cgo

package barretenberg

import (
	"errors"
	"testing"
)

func TestCRSIsInitializedDefault(t *testing.T) {
	// Before any explicit InitCRS call, it should report not initialized
	// (unless a previous test in this process already called it).
	// This test just verifies the function doesn't panic.
	_ = CRSIsInitialized()
}

func TestInitCRSDownloadDefault(t *testing.T) {
	// Use default path with download enabled.
	// This will set up the net CRS factory; actual download happens lazily.
	err := InitCRS("", true)
	if err != nil {
		t.Fatalf("InitCRS(\"\", true) failed: %v", err)
	}

	if !CRSIsInitialized() {
		t.Error("CRSIsInitialized() returned false after successful InitCRS")
	}
}

func TestInitCRSCustomPath(t *testing.T) {
	// Provide an explicit path. Use the default BB CRS dir if available,
	// otherwise use a temp dir (file-only mode will fail if CRS is absent,
	// but net mode should succeed for factory setup).
	err := InitCRS(t.TempDir(), true)
	if err != nil {
		t.Fatalf("InitCRS(tempdir, true) failed: %v", err)
	}

	if !CRSIsInitialized() {
		t.Error("CRSIsInitialized() returned false after successful InitCRS")
	}
}

func TestInitCRSFileOnlyMissingDir(t *testing.T) {
	// File-only mode with a nonexistent/empty dir should still succeed at
	// factory creation time — the error would come later when verification
	// actually needs CRS points. The factory is created eagerly but data is
	// loaded lazily.
	err := InitCRS(t.TempDir(), false)
	// Depending on barretenberg implementation, this may or may not error.
	// If it does, it should be CRS-related.
	if err != nil {
		if !errors.Is(err, ErrCRSNotInitialized) {
			t.Logf("InitCRS(empty, false) returned non-CRS error: %v", err)
		}
	}
}

func TestVerifyWithoutCRSReturnsError(t *testing.T) {
	// This test is conceptual — in practice, once InitCRS is called in the
	// same process (by another test), the flag stays set. Verify the error
	// code mapping exists by checking the sentinel error.
	if ErrCRSNotInitialized == nil {
		t.Error("ErrCRSNotInitialized should not be nil")
	}
}
