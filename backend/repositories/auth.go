package repositories

import (
	"errors"

	dbpkg "lyrics/db"
	"lyrics/models"

	"gorm.io/gorm"
)

var ErrEmailAlreadyExists = errors.New("email already exists")

func EmailExists(email string) (bool, error) {
	if dbpkg.Db == nil {
		return false, errors.New("db not initialized")
	}

	var user models.User
	err := dbpkg.Db.Select("id").Where("email = ?", email).First(&user).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return false, nil
	}
	if err != nil {
		return false, err
	}

	return true, nil
}

func CreateUser(user *models.User) error {
	if dbpkg.Db == nil {
		return errors.New("db not initialized")
	}
	if user == nil {
		return errors.New("user is nil")
	}

	exists, err := EmailExists(user.Email)
	if err != nil {
		return err
	}
	if exists {
		return ErrEmailAlreadyExists
	}

	return dbpkg.Db.Create(user).Error
}
