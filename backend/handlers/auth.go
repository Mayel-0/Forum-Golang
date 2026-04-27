package handlers

import (
	"errors"
	"log"

	"lyrics/models"
	"net/http"

	repositoriespkg "lyrics/repositories"

	"golang.org/x/crypto/bcrypt"
)

func LoginHandle(w http.ResponseWriter, r *http.Request) {
	if tpl == nil {
		http.Error(w, "templates non initialisés", http.StatusInternalServerError)
		return
	}

	if err := tpl.ExecuteTemplate(w, "login.html", nil); err != nil {
		http.Error(w, "Erreur lors du rendu de la page de connexion", http.StatusInternalServerError)
		log.Printf("Erreur template: %v", err)
	}
}

func RegisterHandle(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		if tpl == nil {
			http.Error(w, "templates non initialisés", http.StatusInternalServerError)
			return
		}

		if err := tpl.ExecuteTemplate(w, "register.html", nil); err != nil {
			http.Error(w, "Erreur lors du rendu de la page d'inscription", http.StatusInternalServerError)
			log.Printf("Erreur template: %v", err)
		}
	case http.MethodPost:
		user := models.User{
			Username: r.FormValue("username"),
			Email:    r.FormValue("email"),
		}

		password := r.FormValue("password")
		vPassword := r.FormValue("confirm_password")

		if password == "" {
			http.Error(w, "Mot de passe requis", http.StatusBadRequest)
			return
		}

		if password != vPassword {
			http.Error(w, "Les mots de passe ne correspondent pas", http.StatusBadRequest)
			return
		}

		hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
		if err != nil {
			http.Error(w, "Erreur lors du chiffrement du mot de passe", http.StatusInternalServerError)
			log.Printf("bcrypt error: %v", err)
			return
		}

		user.PasswordHash = string(hash)

		if err := repositoriespkg.CreateUser(&user); err != nil {
			if errors.Is(err, repositoriespkg.ErrEmailAlreadyExists) {
				http.Error(w, "email déjà utilisé", http.StatusConflict)
				return
			}
			http.Error(w, "Erreur lors de la création du compte", http.StatusInternalServerError)
			log.Printf("create user error: %v", err)
			return
		}

		http.Redirect(w, r, "/auth/login", http.StatusSeeOther)

	default:
		http.Error(w, "méthode non autorisée", http.StatusMethodNotAllowed)
	}
}
