package database_test

import (
	"testing"

	"github.com/RAF-SI-2025/EXBanka-3-Infrastructure/internal/database"
)

// --- SeedCurrencies data tests ---

func TestDefaultCurrencies_Count(t *testing.T) {
	if len(database.DefaultCurrencies) != 8 {
		t.Errorf("DefaultCurrencies length = %d, want 8", len(database.DefaultCurrencies))
	}
}

func TestDefaultCurrencies_ContainsExpectedCodes(t *testing.T) {
	expected := []string{"RSD", "EUR", "USD", "CHF", "GBP", "JPY", "CAD", "AUD"}
	codeSet := make(map[string]bool, len(database.DefaultCurrencies))
	for _, c := range database.DefaultCurrencies {
		codeSet[c.Kod] = true
	}
	for _, code := range expected {
		if !codeSet[code] {
			t.Errorf("DefaultCurrencies missing expected code %q", code)
		}
	}
}

func TestDefaultCurrencies_NoDuplicateCodes(t *testing.T) {
	seen := make(map[string]bool)
	for _, c := range database.DefaultCurrencies {
		if seen[c.Kod] {
			t.Errorf("DefaultCurrencies has duplicate code %q", c.Kod)
		}
		seen[c.Kod] = true
	}
}

func TestDefaultCurrencies_AllHaveNaziv(t *testing.T) {
	for _, c := range database.DefaultCurrencies {
		if c.Naziv == "" {
			t.Errorf("Currency %q has empty Naziv", c.Kod)
		}
	}
}

// --- SeedSifreDelatnosti data tests ---

func TestDefaultSifreDelatnosti_Count(t *testing.T) {
	if len(database.DefaultSifreDelatnosti) < 5 {
		t.Errorf("DefaultSifreDelatnosti length = %d, want at least 5", len(database.DefaultSifreDelatnosti))
	}
}

func TestDefaultSifreDelatnosti_ContainsExpectedCodes(t *testing.T) {
	expected := []string{"6419", "6492"}
	codeSet := make(map[string]bool, len(database.DefaultSifreDelatnosti))
	for _, s := range database.DefaultSifreDelatnosti {
		codeSet[s.Sifra] = true
	}
	for _, code := range expected {
		if !codeSet[code] {
			t.Errorf("DefaultSifreDelatnosti missing expected code %q", code)
		}
	}
}

func TestDefaultSifreDelatnosti_NoDuplicateCodes(t *testing.T) {
	seen := make(map[string]bool)
	for _, s := range database.DefaultSifreDelatnosti {
		if seen[s.Sifra] {
			t.Errorf("DefaultSifreDelatnosti has duplicate code %q", s.Sifra)
		}
		seen[s.Sifra] = true
	}
}

func TestDefaultSifreDelatnosti_AllHaveNaziv(t *testing.T) {
	for _, s := range database.DefaultSifreDelatnosti {
		if s.Naziv == "" {
			t.Errorf("SifraDelatnosti %q has empty Naziv", s.Sifra)
		}
	}
}

// --- SeedSifrePlacanja data tests ---

func TestDefaultSifrePlacanja_Count(t *testing.T) {
	if len(database.DefaultSifrePlacanja) < 5 {
		t.Errorf("DefaultSifrePlacanja length = %d, want at least 5", len(database.DefaultSifrePlacanja))
	}
}

func TestDefaultSifrePlacanja_ContainsExpectedCodes(t *testing.T) {
	expected := []string{"221", "222", "289"}
	codeSet := make(map[string]bool, len(database.DefaultSifrePlacanja))
	for _, s := range database.DefaultSifrePlacanja {
		codeSet[s.Sifra] = true
	}
	for _, code := range expected {
		if !codeSet[code] {
			t.Errorf("DefaultSifrePlacanja missing expected code %q", code)
		}
	}
}

func TestDefaultSifrePlacanja_NoDuplicateCodes(t *testing.T) {
	seen := make(map[string]bool)
	for _, s := range database.DefaultSifrePlacanja {
		if seen[s.Sifra] {
			t.Errorf("DefaultSifrePlacanja has duplicate code %q", s.Sifra)
		}
		seen[s.Sifra] = true
	}
}

func TestDefaultSifrePlacanja_AllHaveNaziv(t *testing.T) {
	for _, s := range database.DefaultSifrePlacanja {
		if s.Naziv == "" {
			t.Errorf("SifraPlacanja %q has empty Naziv", s.Sifra)
		}
	}
}
