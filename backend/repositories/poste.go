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

func GetPostBySlug(slug string) (models.Post, error) {
	if dbpkg.Db == nil {
		return models.Post{}, errors.New("db not initialized")
	}

	var post models.Post
	if err := dbpkg.Db.
		Preload("Author").
		Where("slug = ? AND deleted_at IS NULL", slug).
		First(&post).Error; err != nil {
		return models.Post{}, err
	}

	return post, nil
}

func GetAllPostsByCategory(categoryID uuid.UUID) ([]models.Post, error) {
	if dbpkg.Db == nil {
		return nil, errors.New("db not initialized")
	}

	var posts []models.Post
	if err := dbpkg.Db.
		Preload("Author").
		Where("category_id = ? AND deleted_at IS NULL", categoryID).
		Order("created_at DESC").
		Find(&posts).Error; err != nil {
		return nil, err
	}

	return posts, nil
}

func GetAllPostsByCategoryLimit(categoryID uuid.UUID) ([]models.Post, error) {
	if dbpkg.Db == nil {
		return nil, errors.New("db not initialized")
	}

	var posts []models.Post
	if err := dbpkg.Db.
		Preload("Author").
		Where("category_id = ? AND deleted_at IS NULL", categoryID).
		Order("created_at DESC").
		Limit(10).
		Find(&posts).Error; err != nil {
		return nil, err
	}

	return posts, nil
}

func GetAllPostsByCategoryLimitArtistes(categoryID uuid.UUID) ([]models.Post, error) {
	if dbpkg.Db == nil {
		return nil, errors.New("db not initialized")
	}

	var posts []models.Post
	if err := dbpkg.Db.
		Preload("Author").
		Where("category_id = ? AND sub_category = ? AND deleted_at IS NULL", categoryID, models.SubCategoryArtistes).
		Order("created_at DESC").
		Limit(5).
		Find(&posts).Error; err != nil {
		return nil, err
	}

	return posts, nil
}

func GetAllPostsByCategoryLimitConcerts(categoryID uuid.UUID) ([]models.Post, error) {
	if dbpkg.Db == nil {
		return nil, errors.New("db not initialized")
	}

	var posts []models.Post
	if err := dbpkg.Db.
		Preload("Author").
		Where("category_id = ? AND sub_category = ? AND deleted_at IS NULL", categoryID, models.SubCategoryConcerts).
		Order("created_at DESC").
		Limit(5).
		Find(&posts).Error; err != nil {
		return nil, err
	}

	return posts, nil
}

func GetAllPostsByCategoryLimitNouveautes(categoryID uuid.UUID) ([]models.Post, error) {
	if dbpkg.Db == nil {
		return nil, errors.New("db not initialized")
	}

	var posts []models.Post
	if err := dbpkg.Db.
		Preload("Author").
		Where("category_id = ? AND sub_category = ? AND deleted_at IS NULL", categoryID, models.SubCategoryNouveautes).
		Order("created_at DESC").
		Limit(5).
		Find(&posts).Error; err != nil {
		return nil, err
	}

	return posts, nil
}

func GetAllPostsByCategoryAndSub(categoryID uuid.UUID, sub models.PostSubCategory) ([]models.Post, error) {
	if dbpkg.Db == nil {
		return nil, errors.New("db not initialized")
	}
	var posts []models.Post
	if err := dbpkg.Db.
		Preload("Author").
		Where("category_id = ? AND sub_category = ? AND deleted_at IS NULL", categoryID, sub).
		Order("created_at DESC").
		Find(&posts).Error; err != nil {
		return nil, err
	}
	return posts, nil
}

<<<<<<< HEAD
func GetRecentPosts(limit int) ([]models.Post, error) {
	if dbpkg.Db == nil {
		return nil, errors.New("db not initialized")
	}
	var posts []models.Post
	if err := dbpkg.Db.
		Preload("Author").
		Where("deleted_at IS NULL").
		Order("created_at DESC").
		Limit(limit).
		Find(&posts).Error; err != nil {
		return nil, err
	}
	return posts, nil
}

func GetPopularPosts(limit int) ([]models.Post, error) {
	if dbpkg.Db == nil {
		return nil, errors.New("db not initialized")
	}
	var posts []models.Post
	if err := dbpkg.Db.
		Preload("Author").
		Where("deleted_at IS NULL").
		Order("likes_count DESC, comments_count DESC").
		Limit(limit).
		Find(&posts).Error; err != nil {
		return nil, err
	}
	return posts, nil
}

=======
>>>>>>> comment
func CountPostsByCategoryAndSub(categoryID uuid.UUID, subCategory models.PostSubCategory) (int64, error) {
	var count int64
	err := dbpkg.Db.Model(&models.Post{}).Where("category_id = ? AND sub_category = ? AND deleted_at IS NULL", categoryID, subCategory).Count(&count).Error
	return count, err
}
