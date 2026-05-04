package handlers

import (
	"log"
	"lyrics/auth"
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
