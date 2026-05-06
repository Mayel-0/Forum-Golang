package repositories

import (
	"errors"
	dbpkg "lyrics/db"
	"lyrics/models"

	"github.com/google/uuid"
)

func AddComment(comment *models.Comment) error {
	if dbpkg.Db == nil {
		return errors.New("db not initialized")
	}

	return dbpkg.Db.Create(comment).Error
}

func DeleteComment(commentID uuid.UUID) error {
	if dbpkg.Db == nil {
		return errors.New("db not initialized")
	}

	return dbpkg.Db.Delete(&models.Comment{}, commentID).Error
}

func GetCommentsByPostID(postID uuid.UUID) ([]models.Comment, error) {
	if dbpkg.Db == nil {
		return nil, errors.New("db not initialized")
	}

	var comments []models.Comment
	if err := dbpkg.Db.Where("post_id = ?", postID).Order("created_at ASC").Find(&comments).Error; err != nil {
		return nil, err
	}

	return comments, nil
}

func GetCommentByID(commentID uuid.UUID) (*models.Comment, error) {
	if dbpkg.Db == nil {
		return nil, errors.New("db not initialized")
	}

	var comment models.Comment
	if err := dbpkg.Db.Where("id = ?", commentID).First(&comment).Error; err != nil {
		return nil, err
	}

	return &comment, nil
}
