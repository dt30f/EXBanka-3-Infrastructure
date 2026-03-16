package util

import (
	"fmt"
	"math/rand"
	"strconv"
)

// BankCode is the 3-digit bank identifier prefix for all generated account numbers.
const BankCode = "000"

// GenerateAccountNumber generates a random valid 18-digit account number.
// Format: BBB (3-digit bank code) + CCCCCCCCCCCCCC (13-digit account) + KK (2-digit check digits).
// Check digits satisfy: full_number mod 97 == 1 (ISO 7064 MOD 97-10).
func GenerateAccountNumber() string {
	// 13-digit account part in range [1000000000000, 9999999999999]
	accountPart := fmt.Sprintf("%013d", rand.Int63n(9_000_000_000_000)+1_000_000_000_000)
	base := BankCode + accountPart
	return base + checkDigits(base)
}

// ValidateAccountNumber returns true if the number is a valid 18-digit account number.
// Checks: exactly 18 digits, all numeric, and full_number mod 97 == 1.
func ValidateAccountNumber(number string) bool {
	if len(number) != 18 {
		return false
	}
	for _, c := range number {
		if c < '0' || c > '9' {
			return false
		}
	}
	n, err := strconv.ParseUint(number, 10, 64)
	if err != nil {
		return false
	}
	return n%97 == 1
}

// checkDigits computes the 2-digit check suffix for a 16-digit base string.
// kk = 98 - ((base + "00") mod 97), ensuring (base + kk) mod 97 == 1.
func checkDigits(base string) string {
	n, _ := strconv.ParseUint(base+"00", 10, 64)
	kk := 98 - n%97
	return fmt.Sprintf("%02d", kk)
}
