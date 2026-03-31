//go:build cgo
// +build cgo

package barretenberg

import (
	"strings"
	"testing"
)

// TestParseDecimalFieldElementLargeValue verifies that decimal values larger than
// 2^64 are parsed correctly without uint64 overflow.
//
// The old implementation used a uint64 accumulator; any value > 2^64 silently
// wrapped. For example, 2^64+1 = 18446744073709551617 would have been stored as
// 1 rather than the correct 65-bit value in 32-byte big-endian encoding.
func TestParseDecimalFieldElementLargeValue(t *testing.T) {
	// 2^64 + 1 = 18446744073709551617
	// In 32-byte big-endian:  bytes[23]=0x01, bytes[24-30]=0x00, bytes[31]=0x01
	large := "18446744073709551617"
	pi, err := ParsePublicInputsFromStrings([]string{large})
	if err != nil {
		t.Fatalf("unexpected error parsing large decimal input: %v", err)
	}
	if pi.Count() != 1 {
		t.Fatalf("expected 1 element, got %d", pi.Count())
	}
	b := pi.Bytes()
	if len(b) != 32 {
		t.Fatalf("expected 32 bytes, got %d", len(b))
	}
	// byte[23] must be 0x01 — the high word of the 65-bit value.
	// The old uint64 overflow would have left byte[23] == 0 and byte[31] == 1.
	if b[23] != 0x01 {
		t.Errorf("bytes[23] = 0x%02x, expected 0x01 — value was likely truncated by uint64 overflow", b[23])
	}
	if b[31] != 0x01 {
		t.Errorf("bytes[31] = 0x%02x, expected 0x01", b[31])
	}
	// All other bytes must be zero
	for i := 0; i < 23; i++ {
		if b[i] != 0 {
			t.Errorf("bytes[%d] = 0x%02x, expected 0x00", i, b[i])
		}
	}
	for i := 24; i < 31; i++ {
		if b[i] != 0 {
			t.Errorf("bytes[%d] = 0x%02x, expected 0x00", i, b[i])
		}
	}
}

// TestParseDecimalFieldElementMaxUint64 verifies that 2^64-1 (uint64 max) parses correctly.
func TestParseDecimalFieldElementMaxUint64(t *testing.T) {
	maxUint64 := "18446744073709551615" // 2^64 - 1 = 0xffffffffffffffff
	pi, err := ParsePublicInputsFromStrings([]string{maxUint64})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	b := pi.Bytes()
	// In 32-byte big-endian: bytes[0-23] = 0x00, bytes[24-31] = 0xff
	for i := 0; i < 24; i++ {
		if b[i] != 0 {
			t.Errorf("bytes[%d] = 0x%02x, expected 0x00", i, b[i])
		}
	}
	for i := 24; i < 32; i++ {
		if b[i] != 0xff {
			t.Errorf("bytes[%d] = 0x%02x, expected 0xff", i, b[i])
		}
	}
}

// TestParseHexFieldElementShortString verifies that short hex strings (< 64 chars)
// are correctly left-padded to 32 bytes.
//
// The old implementation used fmt.Sprintf("%064s", hexStr) which pads with spaces,
// then manually replaced spaces with zeros. strings.Repeat is now used instead.
func TestParseHexFieldElementShortString(t *testing.T) {
	tests := []struct {
		name      string
		input     string
		wantByte  byte // expected value of the last byte
		wantIndex int  // byte index of the expected non-zero value
	}{
		{
			name:      "single byte 0x2a",
			input:     "0x2a",
			wantByte:  0x2a,
			wantIndex: 31,
		},
		{
			name:      "single byte 0x01",
			input:     "0x01",
			wantByte:  0x01,
			wantIndex: 31,
		},
		{
			name:      "two bytes 0x0102",
			input:     "0x0102",
			wantByte:  0x01,
			wantIndex: 30,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			pi, err := ParsePublicInputsFromStrings([]string{tc.input})
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			b := pi.Bytes()
			if len(b) != 32 {
				t.Fatalf("expected 32 bytes, got %d", len(b))
			}
			if b[tc.wantIndex] != tc.wantByte {
				t.Errorf("bytes[%d] = 0x%02x, expected 0x%02x", tc.wantIndex, b[tc.wantIndex], tc.wantByte)
			}
			// All bytes before the expected non-zero index should be zero
			for i := 0; i < tc.wantIndex; i++ {
				if b[i] != 0 {
					t.Errorf("bytes[%d] = 0x%02x, expected 0x00 (leading zero)", i, b[i])
				}
			}
		})
	}
}

// TestParseDecimalFieldElementInvalid verifies error handling for invalid decimal strings.
func TestParseDecimalFieldElementInvalid(t *testing.T) {
	tests := []struct {
		name  string
		input string
	}{
		{name: "empty", input: ""},
		{name: "non-digit", input: "12a34"},
		{name: "negative", input: "-1"},
		{name: "leading plus", input: "+42"},
		{name: "overflow 33 bytes", input: "115792089237316195423570985008687907853269984665640564039457584007913129639936"}, // 2^256
		{name: "too long 79 significant digits", input: "1000000000000000000000000000000000000000000000000000000000000000000000000000000"}, // 79 significant digits
		{name: "too long absolute", input: strings.Repeat("0", 257)}, // 257 chars — exceeds absolute DoS cap (even though value is 0)
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			_, err := ParsePublicInputsFromStrings([]string{tc.input})
			if err == nil {
				t.Errorf("expected error for input %q, got nil", tc.input)
			}
		})
	}
}

// TestParseDecimalFieldElementLeadingZeros verifies that inputs with leading zeros
// are accepted and parse to the same value as without leading zeros.
func TestParseDecimalFieldElementLeadingZeros(t *testing.T) {
	// "007" should parse identically to "7"
	pi1, err := ParsePublicInputsFromStrings([]string{"007"})
	if err != nil {
		t.Fatalf("unexpected error for leading-zero input: %v", err)
	}
	pi2, err := ParsePublicInputsFromStrings([]string{"7"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	b1 := pi1.Bytes()
	b2 := pi2.Bytes()
	for i := range b1 {
		if b1[i] != b2[i] {
			t.Errorf("bytes[%d]: leading-zero input = 0x%02x, canonical = 0x%02x", i, b1[i], b2[i])
		}
	}

	// 78 leading zeros + "1" = 79 chars total but only 1 significant digit — must be accepted
	padded := strings.Repeat("0", 78) + "1"
	_, err = ParsePublicInputsFromStrings([]string{padded})
	if err != nil {
		t.Errorf("unexpected error for 78-leading-zero input (79 chars, 1 significant digit): %v", err)
	}
}
