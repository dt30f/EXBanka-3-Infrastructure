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
