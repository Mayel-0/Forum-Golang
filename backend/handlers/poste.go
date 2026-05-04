package handlers

import (
	"log"
	"lyrics/auth"
	"lyrics/db"
	"lyrics/models"
	"net/http"

	"github.com/google/uuid"
)

func PosteCreateHandler(w http.ResponseWriter, r *http.Request) {
	UserID, ok := auth.GetUserID(r)
	if !ok {
		http.Error(w, "Utilisateur non authentifié", http.StatusUnauthorized)
		return
	}
	switch r.Method {
	case http.MethodGet:
		render(w, "post-create.html", nil)
	case http.MethodPost:
		userIDParsed, err := uuid.Parse(UserID)
		if err != nil {
			http.Error(w, "ID utilisateur invalide", http.StatusBadRequest)
			return
		}

		Post := models.Post{
			AuthorID: userIDParsed,
			Title:    r.FormValue("title"),
			Body:     r.FormValue("body"),
		}

		if err := db.Db.Create(&Post).Error; err != nil {
			http.Error(w, "Erreur lors de la création du post", http.StatusInternalServerError)
			log.Printf("Erreur création post: %v", err)
			return
		}
		log.Printf("Nouveau post: %+v", Post)
	default:
		http.Error(w, "Méthode non autorisée", http.StatusMethodNotAllowed)
	}
}

func PosteDeleteHandler(w http.ResponseWriter, r *http.Request) {
	UserID, ok := auth.GetUserID(r)
	if !ok {
		http.Error(w, "Utilisateur non authentifié", http.StatusUnauthorized)
		return
	}
	switch r.Method {
	case http.MethodPost:
		log.Println(UserID)
	default:
		http.Error(w, "Méthode non autorisée", http.StatusMethodNotAllowed)
	}
}

func PosteModifierHandle(w http.ResponseWriter, r *http.Request) {
	UserID, ok := auth.GetUserID(r)
	if !ok {
		http.Error(w, "Utilisateur non authentifié", http.StatusUnauthorized)
		return
	}
	switch r.Method {
	case http.MethodPost:
		log.Println(UserID)
	default:
		http.Error(w, "Méthode non autorisée", http.StatusMethodNotAllowed)
	}
}
