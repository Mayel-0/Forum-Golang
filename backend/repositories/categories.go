package repositories

import (
	dbpkg "lyrics/db"
	"lyrics/models"

	"gorm.io/gorm"
)

func GetAllCategories() ([]models.Category, error) {
	var categories []models.Category

	if err := dbpkg.Db.Order("name ASC").Find(&categories).Error; err != nil {
		return nil, err
	}

	return categories, nil
}

func GetCategoryByName(name string) (*models.Category, error) {
	var category models.Category

	err := dbpkg.Db.Where("name = ?", name).First(&category).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, err
	}

	return &category, nil
}
