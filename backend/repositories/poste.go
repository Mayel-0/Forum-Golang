package repositories

import (
	"errors"
	"regexp"
	"strings"

	dbpkg "lyrics/db"
	"lyrics/models"

	"github.com/google/uuid"
)

func slugify(s string) string {
	// Convertir en minuscules
	s = strings.ToLower(s)
	// Remplacer les espaces par des tirets
	s = strings.ReplaceAll(s, " ", "-")
	// Supprimer les caractères spéciaux, garder seulement lettres, chiffres et tirets
	reg := regexp.MustCompile("[^a-z0-9-]+")
	s = reg.ReplaceAllString(s, "")
	// Supprimer les tirets multiples
	reg2 := regexp.MustCompile("-+")
	s = reg2.ReplaceAllString(s, "-")
	// Supprimer les tirets au début et à la fin
	s = strings.Trim(s, "-")
	return s
}

func CreatePoste(Post *models.Post) error {
	if dbpkg.Db == nil {
		return errors.New("db not initialized")
	}
	return dbpkg.Db.Create(&Post).Error
}

func UpdatePoste(Post *models.Post) error {
	if dbpkg.Db == nil {
		return errors.New("db not initialized")
	}
	if Post == nil {
		return errors.New("post is nil")
	}

	return dbpkg.Db.Model(&models.Post{}).
		Where("id = ?", Post.ID).
		Updates(map[string]interface{}{
			"category_id": Post.CategoryID,
			"slug":        Post.Slug,
			"title":       Post.Title,
			"body":        Post.Body,
			"is_pinned":   Post.IsPinned,
			"is_locked":   Post.IsLocked,
		}).Error
}

func DeletePoste(Post *models.Post) error {
	if dbpkg.Db == nil {
		return errors.New("db not initialized")
	}

	return dbpkg.Db.Delete(&Post).Error
}

func UpdatePostSlugFromTitle(Slug string, newTitle string) (string, error) {
	if dbpkg.Db == nil {
		return "", errors.New("db not initialized")
	}

	lastSlashIdx := strings.LastIndex(Slug, "/")
	var prefix string
	if lastSlashIdx != -1 {
		prefix = Slug[:lastSlashIdx+1]
	}

	newSlug := prefix + slugify(newTitle)

	return newSlug, nil
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

func GetAllPostsByCategoryLimit(categoryID uuid.UUID) ([]models.Post, error) {
	if dbpkg.Db == nil {
		return nil, errors.New("db not initialized")
	}

	var posts []models.Post
	if err := dbpkg.Db.Where("category_id = ?", categoryID).Order("name ASC").Limit(10).Find(&posts).Error; err != nil {
		return nil, err
	}

	return posts, nil
}

func GetAllPostsByCategory(categoryID uuid.UUID) ([]models.Post, error) {
	if dbpkg.Db == nil {
		return nil, errors.New("db not initialized")
	}

	var posts []models.Post
	if err := dbpkg.Db.Where("category_id = ?", categoryID).Order("name ASC").Find(&posts).Error; err != nil {
		return nil, err
	}

	return posts, nil
}

func GetCategoryIDBySlug(slug string) (uuid.UUID, error) {
	if dbpkg.Db == nil {
		return uuid.Nil, errors.New("db not initialized")
	}

	var category models.Category
	if err := dbpkg.Db.Where("slug = ?", slug).First(&category).Error; err != nil {
		return uuid.Nil, err
	}

	return category.ID, nil
}

func GetPostByID(postID uuid.UUID) (models.Post, error) {
	if dbpkg.Db == nil {
		return models.Post{}, errors.New("db not initialized")
	}

	var post models.Post
	if err := dbpkg.Db.Where("id = ?", postID).First(&post).Error; err != nil {
		return models.Post{}, err
	}

	return post, nil
}

func GetAuthorByPostID(postID uuid.UUID) (models.User, error) {
	if dbpkg.Db == nil {
		return models.User{}, errors.New("db not initialized")
	}

	var post models.Post
	if err := dbpkg.Db.Where("id = ?", postID).First(&post).Error; err != nil {
		return models.User{}, err
	}

	var author models.User
	if err := dbpkg.Db.Where("id = ?", post.AuthorID).First(&author).Error; err != nil {
		return models.User{}, err
	}

	return author, nil
}
