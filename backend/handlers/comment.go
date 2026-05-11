package handlers

import (
	"encoding/json"
	"lyrics/auth"
	"lyrics/models"
	"net/http"

	repo "lyrics/repositories"

	"github.com/google/uuid"
)

func AddCommentHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodPost:
		UserId, ok := auth.GetUserID(r)
		if !ok {
			http.Error(w, "Utilisateur non authentifié", http.StatusUnauthorized)
			return
		}
		UserIdParsed, err := uuid.Parse(UserId)
		if err != nil {
			http.Error(w, "ID utilisateur invalide", http.StatusBadRequest)
			return
		}

		postID, err := uuid.Parse(r.FormValue("post_id"))
		if err != nil {
			http.Error(w, "ID post invalide", http.StatusBadRequest)
			return
		}

		var parentID *uuid.UUID
		if raw := r.FormValue("parent_id"); raw != "" {
			if pid, err := uuid.Parse(raw); err == nil {
				parentID = &pid
			}
		}

		comment := models.Comment{
			AuthorID: UserIdParsed,
			PostID:   postID,
			Body:     r.FormValue("content"),
			ParentID: parentID,
		}
		if err := repo.AddComment(&comment); err != nil {
			http.Error(w, "Erreur lors de la création du commentaire", http.StatusInternalServerError)
			return
		}

		if r.Header.Get("X-Requested-With") == "fetch" {
			user, _ := auth.GetCurrentUser(r)
			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(map[string]any{
				"id":                comment.ID,
				"body":              comment.Body,
				"created_at":        comment.CreatedAt,
				"author_username":   user.Username,
				"author_avatar_url": user.AvatarURL,
			})
			return
		}

		http.Redirect(w, r, r.FormValue("slug"), http.StatusSeeOther)
	default:
		http.Error(w, "Méthode non autorisée", http.StatusMethodNotAllowed)
	}
}

func DeleteCommentHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodPost:
		UserId, ok := auth.GetUserID(r)
		if !ok {
			http.Error(w, "Utilisateur non authentifié", http.StatusUnauthorized)
			return
		}

		UserIdParsed, err := uuid.Parse(UserId)
		if err != nil {
			http.Error(w, "ID utilisateur invalide", http.StatusBadRequest)
			return
		}

		CommentId, err := uuid.Parse(r.FormValue("comment_id"))
		if err != nil {
			http.Error(w, "ID commentaire invalide", http.StatusBadRequest)
			return
		}

		comment, err := repo.GetCommentByID(CommentId)
		if err != nil {
			http.Error(w, "Commentaire non trouvé", http.StatusNotFound)
			return
		}

		if !auth.VerifyUserRequest(UserIdParsed, comment.AuthorID) {
			http.Error(w, "Utilisateur non autorisé à modifier ce commentaire", http.StatusForbidden)
			return
		}

		if err := repo.DeleteComment(CommentId); err != nil {
			http.Error(w, "Erreur lors de la suppression du commentaire", http.StatusInternalServerError)
			return
		}

	default:
		http.Error(w, "Méthode non autorisée", http.StatusMethodNotAllowed)
	}
}
