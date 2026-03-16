package database

import (
	"fmt"
	"log/slog"

	"github.com/RAF-SI-2025/EXBanka-3-Infrastructure/internal/models"
	"gorm.io/gorm"
)

// DefaultCurrencies contains the 8 currencies seeded at startup.
var DefaultCurrencies = []models.Currency{
	{Kod: "RSD", Naziv: "Srpski dinar", Simbol: "дин", Drzava: "Srbija"},
	{Kod: "EUR", Naziv: "Euro", Simbol: "€", Drzava: "Evropska unija"},
	{Kod: "USD", Naziv: "Američki dolar", Simbol: "$", Drzava: "SAD"},
	{Kod: "CHF", Naziv: "Švajcarski franak", Simbol: "Fr", Drzava: "Švajcarska"},
	{Kod: "GBP", Naziv: "Britanska funta", Simbol: "£", Drzava: "Ujedinjeno Kraljevstvo"},
	{Kod: "JPY", Naziv: "Japanski jen", Simbol: "¥", Drzava: "Japan"},
	{Kod: "CAD", Naziv: "Kanadski dolar", Simbol: "C$", Drzava: "Kanada"},
	{Kod: "AUD", Naziv: "Australijski dolar", Simbol: "A$", Drzava: "Australija"},
}

// DefaultSifreDelatnosti contains standard banking activity codes.
var DefaultSifreDelatnosti = []models.SifraDelatnosti{
	{Sifra: "6410", Naziv: "Novčarsko posredovanje centralne banke"},
	{Sifra: "6419", Naziv: "Ostale novčarske usluge"},
	{Sifra: "6420", Naziv: "Delatnosti holding-kompanija"},
	{Sifra: "6491", Naziv: "Finansijski lizing"},
	{Sifra: "6492", Naziv: "Ostale kreditne delatnosti"},
	{Sifra: "6499", Naziv: "Ostale finansijske usluge"},
	{Sifra: "6511", Naziv: "Životno osiguranje"},
	{Sifra: "6512", Naziv: "Ostale delatnosti osiguranja"},
	{Sifra: "6622", Naziv: "Delatnosti zastupnika i posrednika u osiguranju"},
	{Sifra: "6630", Naziv: "Upravljanje fondovima"},
}

// DefaultSifrePlacanja contains standard payment purpose codes.
var DefaultSifrePlacanja = []models.SifraPlacanja{
	{Sifra: "221", Naziv: "Plaćanje robe"},
	{Sifra: "222", Naziv: "Plaćanje usluga"},
	{Sifra: "240", Naziv: "Plaćanje tekućih transfera"},
	{Sifra: "253", Naziv: "Plaćanje poreza i doprinosa"},
	{Sifra: "254", Naziv: "Plaćanje carina i taksi"},
	{Sifra: "265", Naziv: "Plaćanje po osnovu kredita"},
	{Sifra: "270", Naziv: "Plaćanje zarada"},
	{Sifra: "289", Naziv: "Ostala plaćanja"},
	{Sifra: "290", Naziv: "Interna plaćanja"},
}

// SeedCurrencies inserts the 8 default currencies if they don't already exist.
func SeedCurrencies(db *gorm.DB) error {
	for _, c := range DefaultCurrencies {
		currency := c
		result := db.Where(models.Currency{Kod: currency.Kod}).Assign(models.Currency{
			Naziv:  currency.Naziv,
			Simbol: currency.Simbol,
			Drzava: currency.Drzava,
		}).FirstOrCreate(&currency)
		if result.Error != nil {
			return fmt.Errorf("failed to seed currency %q: %w", currency.Kod, result.Error)
		}
	}
	slog.Info("Currencies seeded", "count", len(DefaultCurrencies))
	return nil
}

// SeedSifreDelatnosti inserts default activity codes if they don't already exist.
func SeedSifreDelatnosti(db *gorm.DB) error {
	for _, s := range DefaultSifreDelatnosti {
		entry := s
		result := db.Where(models.SifraDelatnosti{Sifra: entry.Sifra}).Assign(models.SifraDelatnosti{
			Naziv: entry.Naziv,
		}).FirstOrCreate(&entry)
		if result.Error != nil {
			return fmt.Errorf("failed to seed sifra delatnosti %q: %w", entry.Sifra, result.Error)
		}
	}
	slog.Info("Sifre delatnosti seeded", "count", len(DefaultSifreDelatnosti))
	return nil
}

// DefaultBanka is the bank's own Firma record seeded at startup.
var DefaultBanka = models.Firma{
	Naziv:       "EXBanka",
	MaticniBroj: "11111111",
	PIB:         "111111111",
	Adresa:      "Bulevar Kralja Aleksandra 73, Beograd",
	Telefon:     "+381111234567",
}

// BankaAccountBrojevi maps currency codes to fixed 18-digit bank account numbers.
// Bank-owned accounts use the prefix 000 (bank code) with sequential identifiers.
var BankaAccountBrojevi = map[string]string{
	"RSD": "000000000000000101",
	"EUR": "000000000000000201",
	"USD": "000000000000000301",
	"CHF": "000000000000000401",
	"GBP": "000000000000000501",
	"JPY": "000000000000000601",
	"CAD": "000000000000000701",
	"AUD": "000000000000000801",
}

// SeedBanka creates the bank Firma and one account per currency if they don't exist.
// Must be called after SeedCurrencies.
func SeedBanka(db *gorm.DB) error {
	banka := DefaultBanka
	result := db.Where(models.Firma{MaticniBroj: banka.MaticniBroj}).Assign(models.Firma{
		Naziv:   banka.Naziv,
		PIB:     banka.PIB,
		Adresa:  banka.Adresa,
		Telefon: banka.Telefon,
	}).FirstOrCreate(&banka)
	if result.Error != nil {
		return fmt.Errorf("failed to seed bank firma: %w", result.Error)
	}

	for _, currency := range DefaultCurrencies {
		var cur models.Currency
		if err := db.Where("kod = ?", currency.Kod).First(&cur).Error; err != nil {
			return fmt.Errorf("currency %q not found (run SeedCurrencies first): %w", currency.Kod, err)
		}

		brojRacuna, ok := BankaAccountBrojevi[currency.Kod]
		if !ok {
			return fmt.Errorf("no account number defined for currency %q", currency.Kod)
		}

		account := models.Account{
			BrojRacuna: brojRacuna,
			FirmaID:    &banka.ID,
			CurrencyID: cur.ID,
			Tip:        "tekuci",
			Vrsta:      "poslovni",
			Naziv:      "EXBanka " + currency.Kod + " račun",
			Status:     "aktivan",
		}
		if err := db.Where(models.Account{BrojRacuna: brojRacuna}).Assign(account).FirstOrCreate(&account).Error; err != nil {
			return fmt.Errorf("failed to seed bank account for %q: %w", currency.Kod, err)
		}
	}

	slog.Info("Bank seeded", "firma", banka.Naziv, "accounts", len(DefaultCurrencies))
	return nil
}

// SeedSifrePlacanja inserts default payment purpose codes if they don't already exist.
func SeedSifrePlacanja(db *gorm.DB) error {
	for _, s := range DefaultSifrePlacanja {
		entry := s
		result := db.Where(models.SifraPlacanja{Sifra: entry.Sifra}).Assign(models.SifraPlacanja{
			Naziv: entry.Naziv,
		}).FirstOrCreate(&entry)
		if result.Error != nil {
			return fmt.Errorf("failed to seed sifra placanja %q: %w", entry.Sifra, result.Error)
		}
	}
	slog.Info("Sifre placanja seeded", "count", len(DefaultSifrePlacanja))
	return nil
}
