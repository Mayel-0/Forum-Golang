package handlers

import (
	"log"
	"lyrics/auth"
	"lyrics/models"
	repositoriespkg "lyrics/repositories"
	"net/http"
	"text/template"
)

var tpl *template.Template
var err error

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
	switch r.Method {
	case http.MethodGet:
		data := models.Data{
			User: User,
		}

		if err := tpl.ExecuteTemplate(w, "profile.html", data); err != nil {
			http.Error(w, "Erreur lors du rendu de la page de profil", http.StatusInternalServerError)
			log.Printf("Erreur template: %v", err)
		}
	case http.MethodPost:
		Bio := r.FormValue("bio")
		Email := r.FormValue("email")
		Name := r.FormValue("username")

		bioPtr := &Bio
		updatedUser := &models.User{
			ID:       User.ID,
			Username: Name,
			Email:    Email,
			Bio:      bioPtr,
		}

		if err = repositoriespkg.ModifyUser(updatedUser); err != nil {
			http.Error(w, "Erreur lors de la modification de l'utilisateur", http.StatusInternalServerError)
			log.Printf("Erreur modification utilisateur: %v", err)
			return
		}

		http.Redirect(w, r, "/profile", http.StatusSeeOther)
	default:
		http.Error(w, "méthode non autorisée", http.StatusMethodNotAllowed)
	}
}
