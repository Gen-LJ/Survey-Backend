package category

import "survey-backend/pkg/database"

func SeedCategories() error {
    categories := GetAllCategories()

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
	category := &Category{}
	result := database.DB.Where("name = ?", name).First(&category)
	return category, result.Error
}

func FindByID(id uint) (*Category, error) {
	category := &Category{}
	result := database.DB.First(&category, id)
	return category, result.Error
}
