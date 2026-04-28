package handlers

import (
	"log"
	"lyrics/auth"
	"lyrics/db"
	"lyrics/models"
	"net/http"

	"github.com/google/uuid"
)

/*var tpl *template.Template

func SetTemplates(t *template.Template) {
	tpl = t
} */

func LikeHandlerAdd(w http.ResponseWriter, r *http.Request) {
	UserID, ok := auth.GetUserID(r)
	if !ok {
		http.Error(w, "Utilisateur non authentifié", http.StatusUnauthorized)
		return
	}

	switch r.Method {
	case http.MethodPost:
		PostIDStr := r.FormValue("post_id")
		CommentIDStr := r.FormValue("comment_id")

		userUUID, err := uuid.Parse(UserID)
		if err != nil {
			http.Error(w, "ID utilisateur invalide", http.StatusBadRequest)
			log.Printf("Erreur parse user UUID: %v", err)
			return
		}

		like := models.Likes{
			UserID: userUUID,
		}

		if PostIDStr != "" {
			PostUUID, err := uuid.Parse(PostIDStr)
			if err != nil {
				http.Error(w, "ID post invalide", http.StatusBadRequest)
				log.Printf("Erreur parse post UUID: %v", err)
				return
			}
			like.PostID = &PostUUID
		}

		if CommentIDStr != "" {
			CommentUUID, err := uuid.Parse(CommentIDStr)
			if err != nil {
				http.Error(w, "ID commentaire invalide", http.StatusBadRequest)
				log.Printf("Erreur parse comment UUID: %v", err)
				return
			}
			like.CommentID = &CommentUUID
		}

		if err := db.Db.Create(&like).Error; err != nil {
			http.Error(w, "Erreur lors de la création du like", http.StatusInternalServerError)
			log.Printf("Erreur création like: %v", err)
			return
		}

		w.WriteHeader(http.StatusCreated)

	default:
		http.Error(w, "Méthode non autorisée", http.StatusMethodNotAllowed)
	}
}

func LikeHandlerRm(w http.ResponseWriter, r *http.Request) {
	UserID, ok := auth.GetUserID(r)
	if !ok {
		http.Error(w, "Utilisateur non authentifié", http.StatusUnauthorized)
		return
	}

	switch r.Method {
	case http.MethodPost:
		PostIDStr := r.FormValue("post_id")
		CommentIDStr := r.FormValue("comment_id")

		userUUID, err := uuid.Parse(UserID)
		if err != nil {
			http.Error(w, "ID utilisateur invalide", http.StatusBadRequest)
			log.Printf("Erreur parse user UUID: %v", err)
			return
		}

		query := db.Db.Where("user_id = ?", userUUID)

		if PostIDStr != "" {
			PostUUID, err := uuid.Parse(PostIDStr)
			if err != nil {
				http.Error(w, "ID post invalide", http.StatusBadRequest)
				log.Printf("Erreur parse post UUID: %v", err)
				return
			}
			query = query.Where("post_id = ?", PostUUID)
		}
		if CommentIDStr != "" {
			CommentUUID, err := uuid.Parse(CommentIDStr)
			if err != nil {
				http.Error(w, "ID commentaire invalide", http.StatusBadRequest)
				log.Printf("Erreur parse comment UUID: %v", err)
				return
			}
			query = query.Where("comment_id = ?", CommentUUID)
		}

		if err := query.Delete(&models.Likes{}).Error; err != nil {
			http.Error(w, "Erreur lors de la suppression du like", http.StatusInternalServerError)
			log.Printf("Erreur suppression like: %v", err)
			return
		}

		w.WriteHeader(http.StatusNoContent)

	default:
		http.Error(w, "Méthode non autorisée", http.StatusMethodNotAllowed)
	}
}
