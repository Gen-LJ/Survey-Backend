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
