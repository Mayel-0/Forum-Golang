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
