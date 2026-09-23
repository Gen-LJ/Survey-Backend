package category

import "survey-backend/pkg/database"

func SeedCategories() error {
	categories := GetAllCategories()

	// Skip the per-row FirstOrCreate pass once the table is already populated;
	// against a remote database this is the difference between one round trip
	// on boot and one per category.
	var count int64
	if err := database.DB.Model(&Category{}).Count(&count).Error; err != nil {
		return err
	}
	if count >= int64(len(categories)) {
		return nil
	}

	for i := range categories {
		category := categories[i]
		result := database.DB.Where(Category{Name: category.Name}).FirstOrCreate(&category)
		if result.Error != nil {
			return result.Error
		}
	}
	return nil
}

func FindAll() ([]Category, error) {
	var categories []Category
	result := database.DB.Find(&categories)
	return categories, result.Error
}

func FindByName(name string) (*Category, error) {
	var category Category
	result := database.DB.Where("name = ?", name).First(&category)
	return &category, result.Error
}

func FindByID(id uint) (*Category, error) {
	var category Category
	result := database.DB.First(&category, id)
	return &category, result.Error
}
