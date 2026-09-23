package region

import "survey-backend/pkg/database"

// SeedRegions seeds regions into database
func SeedRegions() error {
	regions := GetAllRegions()

	for i := range regions {
		region := regions[i]
		result := database.DB.Where(Region{Name: region.Name, CountryID: region.CountryID}).FirstOrCreate(&region)
		if result.Error != nil {
			return result.Error
		}
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
