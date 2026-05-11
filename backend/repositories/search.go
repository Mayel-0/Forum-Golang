package repositories

import (
	"log"
	"lyrics/models"
	"math"

	"gorm.io/gorm"
)

func SearchPosts(database *gorm.DB, query string, page, perPage int) ([]models.Post, int64, error) {
	offset := (page - 1) * perPage

	var total int64
	var ids []string

	err := database.Raw(`
		SELECT COUNT(DISTINCT p.id)
		FROM posts p
		WHERE p.deleted_at IS NULL
		AND (
			p.search_vector @@ plainto_tsquery('french', ?)
			OR similarity(p.title, ?) > 0.15
		)
	`, query, query).Scan(&total).Error
	if err != nil {
		log.Println("Erreur count:", err)
		return nil, 0, err
	}
	log.Println("Total trouvé:", total)

	err = database.Raw(`
		SELECT p.id
		FROM posts p
		WHERE p.deleted_at IS NULL
		AND (
			p.search_vector @@ plainto_tsquery('french', ?)
			OR similarity(p.title, ?) > 0.15
		)
		ORDER BY (
			ts_rank(p.search_vector, plainto_tsquery('french', ?)) * 2
			+ similarity(p.title, ?)
		) DESC
		LIMIT ? OFFSET ?
	`, query, query, query, query, perPage, offset).Scan(&ids).Error
	if err != nil {
		log.Println("Erreur ids:", err)
		return nil, 0, err
	}
	log.Println("IDs trouvés:", ids)

	if len(ids) == 0 {
		return []models.Post{}, total, nil
	}

	var posts []models.Post
	err = database.
		Preload("Author").
		Where("id IN ?", ids).
		Find(&posts).Error
	if err != nil {
		log.Println("Erreur posts:", err)
		return nil, 0, err
	}
	log.Println("Posts trouvés:", len(posts))

	return posts, total, nil
}
func TotalPages(total int64, perPage int) int {
	return int(math.Ceil(float64(total) / float64(perPage)))
}
