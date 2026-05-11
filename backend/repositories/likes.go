package repositories

import (
	"errors"
	dbpkg "lyrics/db"
	"lyrics/models"

	"github.com/google/uuid"
)

func LikeAdd(like models.Likes) error {
	if dbpkg.Db == nil {
		return errors.New("db not initialized")
	}

	if like.UserID.String() == "" {
		return errors.New("user ID is empty")
	}

	if like.PostID.String() == "" && like.CommentID.String() == "" {
		return errors.New("post ID and comment ID are both empty")
	}

	return dbpkg.Db.Create(like).Error
}

func LikeRemove(like models.Likes) error {
	if dbpkg.Db == nil {
		return errors.New("db not initialized")
	}

	return dbpkg.Db.Delete(like).Error
}

func GetLikesByPostID(postID string) ([]models.Likes, error) {
	if dbpkg.Db == nil {
		return nil, errors.New("db not initialized")
	}

	var likes []models.Likes
	if err := dbpkg.Db.Where("post_id = ?", postID).Find(&likes).Error; err != nil {
		return nil, err
	}

	return likes, nil
}

func HasUserLikedPost(userID, postID uuid.UUID) bool {
	if dbpkg.Db == nil {
		return false
	}
	var count int64
	dbpkg.Db.Model(&models.Likes{}).Where("user_id = ? AND post_id = ?", userID, postID).Count(&count)
	return count > 0
}

func GetUserLikedCommentIDs(userID uuid.UUID, postID uuid.UUID) []string {
	if dbpkg.Db == nil {
		return nil
	}
	var likes []models.Likes
	dbpkg.Db.Where("user_id = ? AND comment_id IS NOT NULL AND comment_id IN (SELECT id FROM comments WHERE post_id = ?)", userID, postID).Find(&likes)
	ids := make([]string, 0, len(likes))
	for _, l := range likes {
		if l.CommentID != nil {
			ids = append(ids, l.CommentID.String())
		}
	}
	return ids
}
