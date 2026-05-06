package repositories

import (
	"bytes"
	"errors"
	"io"
	"log"
	"net/http"
	"os"

	dbpkg "lyrics/db"
	"lyrics/models"

	"gorm.io/gorm"
)

var ErrEmailAlreadyExists = errors.New("email already exists")
var ErrUserNotFound = errors.New("user not found")

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

func FindUserByEmail(email string) (*models.User, error) {
	if dbpkg.Db == nil {
		return nil, errors.New("db not initialized")
	}

	var user models.User
	err := dbpkg.Db.Where("email = ?", email).First(&user).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, ErrUserNotFound
	}
	if err != nil {
		return nil, err
	}

	return &user, nil
}

func ModifyUser(user *models.User) error {
	if dbpkg.Db == nil {
		return errors.New("db not initialized")
	}
	if user == nil {
		return errors.New("user is nil")
	}

	return dbpkg.Db.Model(&models.User{}).
		Where("id = ?", user.ID).
		Updates(map[string]interface{}{
			"username": user.Username,
			"email":    user.Email,
			"bio":      user.Bio,
		}).Error
}

func GetTopUsers(limit int) ([]models.User, error) {
	if dbpkg.Db == nil {
		return nil, errors.New("db not initialized")
	}

	var users []models.User
	err := dbpkg.Db.Limit(limit).Find(&users).Error

	return users, err
}

func UploadAvatar(userID string, ext string, fileBytes []byte) (string, error) {
	supabaseURL := os.Getenv("SUPABASE_URL")
	supabaseKey := os.Getenv("SUPABASE_SERVICE_KEY")

	log.Printf("SUPABASE_URL: %s", supabaseURL)
	log.Printf("SUPABASE_SERVICE_KEY présent: %v", supabaseKey != "")

	fileName := userID + ext
	uploadURL := supabaseURL + "/storage/v1/object/avatars/" + fileName

	log.Printf("Upload URL: %s", uploadURL)

	req, err := http.NewRequest("POST", uploadURL, bytes.NewReader(fileBytes))
	if err != nil {
		log.Printf("Erreur création requête: %v", err)
		return "", err
	}

	req.Header.Set("Authorization", "Bearer "+supabaseKey)
	req.Header.Set("Content-Type", "image/"+ext[1:])
	req.Header.Set("x-upsert", "true")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		log.Printf("Erreur envoi requête: %v", err)
		return "", err
	}
	defer resp.Body.Close()

	// Lire le body de la réponse pour voir l'erreur Supabase
	body, _ := io.ReadAll(resp.Body)
	log.Printf("Status Supabase: %d", resp.StatusCode)
	log.Printf("Réponse Supabase: %s", string(body))

	if resp.StatusCode != http.StatusOK {
		return "", errors.New("upload supabase échoué")
	}

	publicURL := supabaseURL + "/storage/v1/object/public/avatars/" + fileName
	return publicURL, nil
}

func UpdateUserAvatar(userID string, avatarURL string) error {
	return dbpkg.Db.Model(&models.User{}).
		Where("id = ?", userID).
		Update("avatar_url", avatarURL).Error
}

func GetUserByID(userID string) (models.User, error) {
	if dbpkg.Db == nil {
		return models.User{}, errors.New("db not initialized")
	}

	var user models.User
	err := dbpkg.Db.Where("id = ?", userID).First(&user).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return models.User{}, ErrUserNotFound
	}
	if err != nil {
		return models.User{}, err
	}

	return user, nil
}
