package handlers

import (
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

		ParentID := uuid.MustParse(r.FormValue("parent_id"))
		if ParentID == uuid.Nil {
			ParentID = uuid.Nil
		}

		Slug := r.FormValue("slug")

		comment := models.Comment{
			AuthorID: UserIdParsed,
			PostID:   uuid.MustParse(r.FormValue("post_id")),
			Body:     r.FormValue("content"),
			ParentID: &ParentID,
		}
		if err := repo.AddComment(&comment); err != nil {
			http.Error(w, "Erreur lors de la création du commentaire", http.StatusInternalServerError)
			return
		}
		http.Redirect(w, r, Slug, http.StatusSeeOther)
	default:
		http.Error(w, "Méthode non autorisée", http.StatusMethodNotAllowed)
		return
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

		if !auth.VerifyUserRequest(UserIdParsed, post.AuthorID) {
			http.Error(w, "Utilisateur non autorisé à modifier ce post", http.StatusForbidden)
			return
		}

	default:
		http.Error(w, "Méthode non autorisée", http.StatusMethodNotAllowed)
		return
	}
}
