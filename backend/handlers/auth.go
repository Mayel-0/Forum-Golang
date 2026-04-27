package handlers

import (
	"crypto/rand"
	"encoding/base64"
	"errors"
	"log"
	"time"

	authpkg "lyrics/auth"
	"lyrics/models"
	"net/http"

	repositoriespkg "lyrics/repositories"

	"golang.org/x/crypto/bcrypt"
)

func LoginHandle(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		if tpl == nil {
			http.Error(w, "templates non initialisés", http.StatusInternalServerError)
			return
		}

		if err := tpl.ExecuteTemplate(w, "login.html", nil); err != nil {
			http.Error(w, "Erreur lors du rendu de la page de connexion", http.StatusInternalServerError)
			log.Printf("Erreur template: %v", err)
		}
	case http.MethodPost:
		email := r.FormValue("email")
		password := r.FormValue("password")

		if email == "" || password == "" {
			http.Error(w, "email et mot de passe requis", http.StatusBadRequest)
			return
		}

		user, err := repositoriespkg.FindUserByEmail(email)
		if err != nil {
			if errors.Is(err, repositoriespkg.ErrUserNotFound) {
				http.Error(w, "identifiants invalides", http.StatusUnauthorized)
				return
			}
			http.Error(w, "erreur serveur", http.StatusInternalServerError)
			log.Printf("find user by email error: %v", err)
			return
		}

		if err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(password)); err != nil {
			http.Error(w, "identifiants invalides", http.StatusUnauthorized)
			return
		}

		token, err := generateSessionToken()
		if err != nil {
			http.Error(w, "erreur serveur", http.StatusInternalServerError)
			log.Printf("generate token error: %v", err)
			return
		}

		expiresAt := time.Now().Add(24 * time.Hour)
		authpkg.StoreSession(token, models.Session{
			Userid: user.ID.String(),
			Expiry: expiresAt,
		})

		http.SetCookie(w, &http.Cookie{
			Name:     "session_token",
			Value:    token,
			Path:     "/",
			Expires:  expiresAt,
			HttpOnly: true,
			Secure:   r.TLS != nil,
			SameSite: http.SameSiteLaxMode,
		})

		http.Redirect(w, r, "/", http.StatusSeeOther)

	default:
		http.Error(w, "méthode non autorisée", http.StatusMethodNotAllowed)
	}
}

func generateSessionToken() (string, error) {
	b := make([]byte, 32)
	if _, err := rand.Read(b); err != nil {
		return "", err
	}
	return base64.RawURLEncoding.EncodeToString(b), nil
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
