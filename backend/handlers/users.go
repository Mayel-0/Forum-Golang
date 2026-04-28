package handlers

import (
	"log"
	"lyrics/auth"
	"lyrics/models"
	"net/http"
	"text/template"
)

var tpl *template.Template

func SetTemplates(t *template.Template) {
	tpl = t
}

func AcceuilHandle(w http.ResponseWriter, r *http.Request) {
	if tpl == nil {
		http.Error(w, "templates non initialisés", http.StatusInternalServerError)
		return
	}

	if err := tpl.ExecuteTemplate(w, "accueil.html", nil); err != nil {
		http.Error(w, "Erreur lors du rendu de la page d'accueil", http.StatusInternalServerError)
		log.Printf("Erreur template: %v", err)
	}
}

func ForumIndexHandle(w http.ResponseWriter, r *http.Request) {
	if tpl == nil {
		http.Error(w, "templates non initialisés", http.StatusInternalServerError)
		return
	}

	if err := tpl.ExecuteTemplate(w, "index.html", nil); err != nil {
		http.Error(w, "Erreur lors du rendu de la page du forum", http.StatusInternalServerError)
		log.Printf("Erreur template: %v", err)
	}
}

func ProfileHandle(w http.ResponseWriter, r *http.Request) {
	if tpl == nil {
		http.Error(w, "templates non initialisés", http.StatusInternalServerError)
		return
	}

	UserID, ok := auth.GetUserID(r)
	if !ok {
		http.Error(w, "Utilisateur non authentifié", http.StatusUnauthorized)
		return
	}
	User, err := auth.GetUserByID(UserID)
	if err != nil {
		http.Error(w, "Erreur lors de la récupération de l'utilisateur", http.StatusInternalServerError)
		log.Printf("Erreur récupération utilisateur: %v", err)
		return
	}

	data := models.Data{
		User: User,
	}

	if err := tpl.ExecuteTemplate(w, "profile.html", data); err != nil {
		http.Error(w, "Erreur lors du rendu de la page de profil", http.StatusInternalServerError)
		log.Printf("Erreur template: %v", err)
	}
}
