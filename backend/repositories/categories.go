package repositories

import (
	"errors"
	dbpkg "lyrics/db"
	"lyrics/models"

	"github.com/google/uuid"
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

func GetCategoryIDByName(categoryName string) (uuid.UUID, error) {
	if dbpkg.Db == nil {
		return uuid.Nil, errors.New("db not initialized")
	}

	var category models.Category
	if err := dbpkg.Db.Where("name = ?", categoryName).First(&category).Error; err != nil {
		return uuid.Nil, err
	}

	return category.ID, nil
}

func GetCategoryIDBySlug(slug string) (uuid.UUID, error) {
	var category models.Category
	// Cherche avec le slash préfixé comme stocké en DB
	if err := dbpkg.Db.Where("slug = ?", "/"+slug).First(&category).Error; err != nil {
		return uuid.Nil, err
	}
	return category.ID, nil
}
