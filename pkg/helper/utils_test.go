package helper

import (
	"testing"
)

func TestValidateNIPFromStr(t *testing.T) {
	tests := []struct {
		nip    string
		result bool
	}{
		{"303120200982140", true},       // Valid NIP
		{"615320220100001", false},      // Invalid gender digit
		{"615219990213000", false},      // Invalid year
		{"615220121302013", false},      // Invalid month
		{"61522012010201334222", false}, // Invalid random digits
	}

	for _, test := range tests {
		t.Run(test.nip, func(t *testing.T) {
			result := ValidateNIPFromStr(test.nip)
			if result != test.result {
				t.Errorf("expected %v, got %v", test.result, result)
			}
		})
	}
}
