package repositories

import (
	"errors"
	dbpkg "lyrics/db"
	"lyrics/models"
)

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
