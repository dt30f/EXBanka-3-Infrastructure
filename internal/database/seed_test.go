package database_test

import (
	"testing"

	"github.com/RAF-SI-2025/EXBanka-3-Infrastructure/internal/database"
)

// --- SeedBanka data tests ---

func TestDefaultBanka_Naziv(t *testing.T) {
	if database.DefaultBanka.Naziv != "EXBanka" {
		t.Errorf("DefaultBanka.Naziv = %q, want %q", database.DefaultBanka.Naziv, "EXBanka")
	}
}

func TestDefaultBanka_MaticniBroj_NotEmpty(t *testing.T) {
	if database.DefaultBanka.MaticniBroj == "" {
		t.Error("DefaultBanka.MaticniBroj should not be empty")
	}
}

func TestDefaultBanka_PIB_NotEmpty(t *testing.T) {
	if database.DefaultBanka.PIB == "" {
		t.Error("DefaultBanka.PIB should not be empty")
	}
}

func TestDefaultBanka_MaticniBroj_8Digits(t *testing.T) {
	mb := database.DefaultBanka.MaticniBroj
	if len(mb) != 8 {
		t.Errorf("DefaultBanka.MaticniBroj length = %d, want 8", len(mb))
	}
	for _, c := range mb {
		if c < '0' || c > '9' {
			t.Errorf("DefaultBanka.MaticniBroj contains non-digit: %q", string(c))
		}
	}
}

func TestDefaultBanka_PIB_9Digits(t *testing.T) {
	pib := database.DefaultBanka.PIB
	if len(pib) != 9 {
		t.Errorf("DefaultBanka.PIB length = %d, want 9", len(pib))
	}
	for _, c := range pib {
		if c < '0' || c > '9' {
			t.Errorf("DefaultBanka.PIB contains non-digit: %q", string(c))
		}
	}
}

func TestBankaAccountBrojevi_Count(t *testing.T) {
	if len(database.BankaAccountBrojevi) != 8 {
		t.Errorf("BankaAccountBrojevi length = %d, want 8", len(database.BankaAccountBrojevi))
	}
}

func TestBankaAccountBrojevi_CoversAllCurrencies(t *testing.T) {
	for _, c := range database.DefaultCurrencies {
		if _, ok := database.BankaAccountBrojevi[c.Kod]; !ok {
			t.Errorf("BankaAccountBrojevi missing entry for currency %q", c.Kod)
		}
	}
}

func TestBankaAccountBrojevi_AllSize18(t *testing.T) {
	for kod, brj := range database.BankaAccountBrojevi {
		if len(brj) != 18 {
			t.Errorf("BankaAccountBrojevi[%q] length = %d, want 18", kod, len(brj))
		}
	}
}

func TestBankaAccountBrojevi_AllDigits(t *testing.T) {
	for kod, brj := range database.BankaAccountBrojevi {
		for _, c := range brj {
			if c < '0' || c > '9' {
				t.Errorf("BankaAccountBrojevi[%q] contains non-digit: %q", kod, string(c))
				break
			}
		}
	}
}

func TestBankaAccountBrojevi_Unique(t *testing.T) {
	seen := make(map[string]string)
	for kod, brj := range database.BankaAccountBrojevi {
		if prev, exists := seen[brj]; exists {
			t.Errorf("BankaAccountBrojevi duplicate account number %q for currencies %q and %q", brj, prev, kod)
		}
		seen[brj] = kod
	}
}

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
