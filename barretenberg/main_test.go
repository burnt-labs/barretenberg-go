//go:build cgo
// +build cgo

package barretenberg

import (
	"os"
	"testing"
)

func TestMain(m *testing.M) {
	// Initialise the CRS before running any verification tests.
	// Download mode is enabled so the CRS factory is ready; actual data
	// is fetched lazily on the first verification that needs it.
	if err := InitCRS("", true); err != nil {
		// If CRS init fails, still run tests — unit tests that don't
		// verify proofs will pass; integration tests will skip/fail.
		os.Stderr.WriteString("warning: InitCRS failed: " + err.Error() + "\n")
	}
	os.Exit(m.Run())
}
