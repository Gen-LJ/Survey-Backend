package country

import "survey-backend/pkg/database"

func SeedCountries() error {
	countries := GetAllCountries()

	// Nearly 200 rows: without this guard every restart replays one query per
	// country before the server can accept traffic.
	var count int64
	if err := database.DB.Model(&Country{}).Count(&count).Error; err != nil {
		return err
	}
	if count >= int64(len(countries)) {
		return nil
	}

	for i := range countries {
		country := countries[i]
		result := database.DB.Where(Country{Name: country.Name}).FirstOrCreate(&country)
		if result.Error != nil {
			return result.Error
		}
	}
	return nil
}

func FindAll() ([]Country, error) {
	var countries []Country
	result := database.DB.Find(&countries)
	return countries, result.Error
}

func FindActive() ([]Country, error) {
	var countries []Country
	result := database.DB.Where("active = ?", true).Find(&countries)
	return countries, result.Error
}

func ToggleActiveByID(id uint, active bool) error {
	result := database.DB.Model(&Country{}).Where("id = ?", id).Update("active", active)
	return result.Error
}

func BulkUpdateActive(ids []uint, active bool) error {
	result := database.DB.Model(&Country{}).Where("id IN ?", ids).Update("active", active)
	return result.Error
}

// FindFiltered returns countries, optionally narrowed to active or inactive.
// A nil filter means every country.
func FindFiltered(active *bool) ([]Country, error) {
	var countries []Country

	query := database.DB.Order("name asc")
	if active != nil {
		query = query.Where("active = ?", *active)
	}

	result := query.Find(&countries)
	return countries, result.Error
}

// FindByCodes resolves ISO codes to countries in one query, keyed by code.
// Seeding uses this instead of hardcoded ids: ids are assigned by the database
// at insert time and differ between engines.
func FindByCodes(codes []string) (map[string]Country, error) {
	found := make(map[string]Country, len(codes))
	if len(codes) == 0 {
		return found, nil
	}

	var countries []Country
	if err := database.DB.Where("code IN ?", codes).Find(&countries).Error; err != nil {
		return nil, err
	}

	for _, c := range countries {
		found[c.Code] = c
	}

	return found, nil
}
