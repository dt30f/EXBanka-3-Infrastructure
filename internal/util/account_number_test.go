package util_test

import (
	"strings"
	"testing"

	"github.com/RAF-SI-2025/EXBanka-3-Infrastructure/internal/util"
)

// --- GenerateAccountNumber tests ---

func TestGenerateAccountNumber_Returns18Digits(t *testing.T) {
	n := util.GenerateAccountNumber()
	if len(n) != 18 {
		t.Errorf("GenerateAccountNumber() length = %d, want 18", len(n))
	}
}

func TestGenerateAccountNumber_AllDigits(t *testing.T) {
	n := util.GenerateAccountNumber()
	for _, c := range n {
		if c < '0' || c > '9' {
			t.Errorf("GenerateAccountNumber() contains non-digit: %q in %q", string(c), n)
		}
	}
}

func TestGenerateAccountNumber_IsValid(t *testing.T) {
	for i := 0; i < 10; i++ {
		n := util.GenerateAccountNumber()
		if !util.ValidateAccountNumber(n) {
			t.Errorf("GenerateAccountNumber() produced invalid number: %q", n)
		}
	}
}

func TestGenerateAccountNumber_Uniqueness(t *testing.T) {
	seen := make(map[string]bool, 1000)
	for i := 0; i < 1000; i++ {
		n := util.GenerateAccountNumber()
		if seen[n] {
			t.Errorf("GenerateAccountNumber() produced duplicate: %q", n)
		}
		seen[n] = true
	}
}

// --- ValidateAccountNumber tests ---

func TestValidateAccountNumber_AcceptsValid(t *testing.T) {
	// Generate a valid number, then validate it
	valid := util.GenerateAccountNumber()
	if !util.ValidateAccountNumber(valid) {
		t.Errorf("ValidateAccountNumber(%q) = false, want true", valid)
	}
}

func TestValidateAccountNumber_RejectsWrongLength(t *testing.T) {
	cases := []string{
		"",
		"123",
		"12345678901234567",  // 17 digits
		"1234567890123456789", // 19 digits
	}
	for _, c := range cases {
		if util.ValidateAccountNumber(c) {
			t.Errorf("ValidateAccountNumber(%q) = true, want false (wrong length)", c)
		}
	}
}

func TestValidateAccountNumber_RejectsNonNumeric(t *testing.T) {
	cases := []string{
		"00000000000000000A", // letter at end
		"abcdefghijklmnopqr", // all letters
		"00000000000000000 ", // space
		"000000000000000-01", // dash
	}
	for _, c := range cases {
		if util.ValidateAccountNumber(c) {
			t.Errorf("ValidateAccountNumber(%q) = true, want false (non-numeric)", c)
		}
	}
}

func TestValidateAccountNumber_RejectsInvalidCheckDigits(t *testing.T) {
	// Generate a valid number and corrupt its check digits
	valid := util.GenerateAccountNumber()
	corrupted := valid[:16] + "00"
	// Only reject if "00" doesn't happen to be the correct check digits
	if corrupted != valid {
		if util.ValidateAccountNumber(corrupted) {
			t.Errorf("ValidateAccountNumber(%q) = true for corrupted check digits, want false", corrupted)
		}
	}
}

func TestValidateAccountNumber_RejectsFlippedDigit(t *testing.T) {
	valid := util.GenerateAccountNumber()
	// Flip a digit in the middle
	runes := []rune(valid)
	if runes[8] == '9' {
		runes[8] = '0'
	} else {
		runes[8] = runes[8] + 1
	}
	corrupted := string(runes)
	if strings.Contains(valid, corrupted) {
		t.Skip("corrupted number accidentally equals original")
	}
	if util.ValidateAccountNumber(corrupted) {
		t.Errorf("ValidateAccountNumber(%q) = true for digit-flipped number, want false", corrupted)
	}
}
