package handlers

import (
	"lyrics/auth"
	"net/http"
)

func PosteCreateHandler(w http.ResponseWriter, r *http.Request) {
	UserID, ok := auth.GetUserID(r)
	if !ok {
		http.Error(w, "Utilisateur non authentifié", http.StatusUnauthorized)
		return
	}
	switch r.Method {
	case http.MethodPost:
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
	default:
		http.Error(w, "Méthode non autorisée", http.StatusMethodNotAllowed)
	}
}
