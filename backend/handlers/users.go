package handlers

import (
	"log"
	"lyrics/auth"
	"lyrics/models"
	repositoriespkg "lyrics/repositories"
	"net/http"
	"text/template"
)

var templates map[string]*template.Template

func SetTemplates(t map[string]*template.Template) {
	templates = t
}

func render(w http.ResponseWriter, name string, data any) {
	t, ok := templates[name]
	if !ok {
		http.Error(w, "template introuvable: "+name, http.StatusInternalServerError)
		return
	}
	if err := t.ExecuteTemplate(w, "template.html", data); err != nil {
		http.Error(w, "Erreur lors du rendu", http.StatusInternalServerError)
		log.Printf("Erreur template %s: %v", name, err)
	}
}

func ForumIndexHandle(w http.ResponseWriter, r *http.Request) {
	topUsers, err := repositoriespkg.GetTopUsers(3)
	if err != nil {
		http.Error(w, "Erreur récupération des utilisateurs", http.StatusInternalServerError)
		return
	}
	data := models.Data{TopUsers: topUsers}
	if user, ok := auth.GetCurrentUser(r); ok {
		data.User = user
	}
	render(w, "index.html", data)
}

func ProfileHandle(w http.ResponseWriter, r *http.Request) {
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
		render(w, "profile.html", models.Data{User: User})
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
