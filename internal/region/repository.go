package region

import (
	"fmt"
	"log"

	"survey-backend/internal/country"
	"survey-backend/pkg/database"
)

// SeedRegions creates the seed regions and repairs any whose country link is
// wrong. Each region is matched by its own code and re-pointed at the country
// that the CountryCode resolves to now, so a database seeded against different
// country ids heals itself on the next boot.
func SeedRegions() error {
	seeds := GetAllRegions()

	codes := make([]string, 0, 4)
	seen := map[string]bool{}
	for _, s := range seeds {
		if !seen[s.CountryCode] {
			seen[s.CountryCode] = true
			codes = append(codes, s.CountryCode)
		}
	}

	countries, err := country.FindByCodes(codes)
	if err != nil {
		return fmt.Errorf("failed to resolve countries for region seeding: %w", err)
	}
	for _, code := range codes {
		if _, ok := countries[code]; !ok {
			return fmt.Errorf("cannot seed regions: no country with code %s", code)
		}
	}

	// One read of the existing rows, so a correct database writes nothing.
	var existing []Region
	if err := database.DB.Find(&existing).Error; err != nil {
		return err
	}
	byCode := make(map[string]Region, len(existing))
	for _, r := range existing {
		byCode[r.Code] = r
	}

	created, repaired := 0, 0
	for _, seed := range seeds {
		want := countries[seed.CountryCode]

		current, found := byCode[seed.Code]
		if !found {
			row := Region{Name: seed.Name, Code: seed.Code, CountryID: want.ID, Active: true}
			if err := database.DB.Create(&row).Error; err != nil {
				return err
			}
			created++
			continue
		}

		if current.CountryID == want.ID && current.Name == seed.Name {
			continue
		}

		if err := database.DB.Model(&Region{}).
			Where("id = ?", current.ID).
			Updates(map[string]any{"name": seed.Name, "country_id": want.ID}).Error; err != nil {
			return err
		}
		repaired++
	}

	if created > 0 || repaired > 0 {
		log.Printf("regions seeded: %d created, %d re-pointed at the right country", created, repaired)
	}

	return nil
}

// FindByCountryID finds regions by country ID
func FindByCountryID(countryID uint) ([]Region, error) {
	var regions []Region
	result := database.DB.Where("country_id = ?", countryID).Find(&regions)
	return regions, result.Error
}

// FindActiveByCountryID finds active regions by country ID
func FindActiveByCountryID(countryID uint) ([]Region, error) {
	var regions []Region
	result := database.DB.Where("country_id = ? AND active = ?", countryID, true).Find(&regions)
	return regions, result.Error
}
