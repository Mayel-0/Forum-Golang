package repositories

import (
	"errors"
	dbpkg "lyrics/db"
	"lyrics/models"
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

func GetLikesByPostID(postID string) ([]models.Likes, int, error) {
	if dbpkg.Db == nil {
		return nil, 0, errors.New("db not initialized")
	}

	var likes []models.Likes
	if err := dbpkg.Db.Where("post_id = ?", postID).Find(&likes).Error; err != nil {
		return nil, 0, err
	}

	return likes, len(likes), nil
}

func GetLikesByCommentID(commentID string) ([]models.Likes, int, error) {
	if dbpkg.Db == nil {
		return nil, 0, errors.New("db not initialized")
	}

	var likes []models.Likes
	if err := dbpkg.Db.Where("comment_id = ?", commentID).Find(&likes).Error; err != nil {
		return nil, 0, err
	}

	return likes, len(likes), nil
}
