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

	"github.com/google/uuid"
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

func GetUsersWithStats(query *gorm.DB) ([]models.User, error) {
	var users []models.User
	if err := query.
		Select(`users.*,
            COUNT(DISTINCT l.id) AS likes_count,
            COUNT(DISTINCT p.id) AS posts_count,
            COUNT(DISTINCT c.id) AS comments_count`).
		Joins("LEFT JOIN likes l ON l.user_id = users.id").
		Joins("LEFT JOIN posts p ON p.author_id = users.id AND p.deleted_at IS NULL").
		Joins("LEFT JOIN comments c ON c.author_id = users.id AND c.deleted_at IS NULL").
		Group("users.id").
		Find(&users).Error; err != nil {
		return nil, err
	}
	return users, nil
}

func GetTopUsers(limit int) ([]models.User, error) {
	if dbpkg.Db == nil {
		return nil, errors.New("db not initialized")
	}

	query := dbpkg.Db.Model(&models.User{}).
		Order("likes_count DESC").
		Limit(limit)
	return GetUsersWithStats(query)
}

func GetUserByUsername(username string) (models.User, error) {
	if dbpkg.Db == nil {
		return models.User{}, errors.New("db not initialized")
	}
	query := dbpkg.Db.Model(&models.User{}).Where("users.username = ?", username)
	users, err := GetUsersWithStats(query)
	if err != nil {
		return models.User{}, err
	}
	if len(users) == 0 {
		return models.User{}, ErrUserNotFound
	}
	return users[0], nil
}

func GetUserByID(userID uuid.UUID) (models.User, error) {
	if dbpkg.Db == nil {
		return models.User{}, errors.New("db not initialized")
	}

	query := dbpkg.Db.Model(&models.User{}).
		Where("users.id = ?", userID)
	users, err := GetUsersWithStats(query)
	if err != nil {
		return models.User{}, err
	}
	if len(users) == 0 {
		return models.User{}, ErrUserNotFound
	}
	return users[0], nil
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
	if dbpkg.Db == nil {
		return errors.New("db not initialized")
	}
	return dbpkg.Db.Model(&models.User{}).
		Where("id = ?", userID).
		Update("avatar_url", avatarURL).Error
}

func UploadBanniere(userID string, ext string, fileBytes []byte) (string, error) {
	supabaseURL := os.Getenv("SUPABASE_URL")
	supabaseKey := os.Getenv("SUPABASE_SERVICE_KEY")

	log.Printf("SUPABASE_URL: %s", supabaseURL)
	log.Printf("SUPABASE_SERVICE_KEY présent: %v", supabaseKey != "")

	fileName := userID + ext
	uploadURL := supabaseURL + "/storage/v1/object/bannieres/" + fileName

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

	body, _ := io.ReadAll(resp.Body)
	log.Printf("Status Supabase: %d", resp.StatusCode)
	log.Printf("Réponse Supabase: %s", string(body))

	if resp.StatusCode != http.StatusOK {
		return "", errors.New("upload supabase échoué")
	}

	publicURL := supabaseURL + "/storage/v1/object/public/bannieres/" + fileName
	return publicURL, nil
}

func UpdateUserBanniere(userID string, banniereURL string) error {
	if dbpkg.Db == nil {
		return errors.New("db not initialized")
	}
	return dbpkg.Db.Model(&models.User{}).
		Where("id = ?", userID).
		Update("baniere_url", banniereURL).Error
}
