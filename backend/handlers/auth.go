package handlers

import (
	"errors"
	"net/http"

	authpkg "lyrics/auth"
	"lyrics/models"
	repositoriespkg "lyrics/repositories"

	"golang.org/x/crypto/bcrypt"
)

func LoginHandle(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		render(w, r, "login.html", nil)
	case http.MethodPost:
		email := r.FormValue("email")
		password := r.FormValue("password")

		if email == "" || password == "" {
			data := models.Data{
				MessagesError: "email et mot de passe requis",
			}
			render(w, r, "login.html", data)
			return
		}

		user, err := repositoriespkg.FindUserByEmail(email)
		if err != nil {
			if errors.Is(err, repositoriespkg.ErrUserNotFound) {
				data := models.Data{
					MessagesError: "identifiants invalides",
				}
				render(w, r, "login.html", data)
				return
			}
			data := models.Data{
				CodeError:     500,
				MessagesError: "erreur serveur",
			}
			Error(w, r, &data)
			return
		}

		if err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(password)); err != nil {
			data := models.Data{
				MessagesError: "identifiants invalides",
			}
			render(w, r, "login.html", data)
			return
		}

		if err := authpkg.SetSession(w, r, user.ID.String()); err != nil {
			data := models.Data{
				CodeError:     500,
				MessagesError: "erreur serveur",
			}
			Error(w, r, &data)
			return
		}

		http.Redirect(w, r, "/", http.StatusSeeOther)

	default:
		data := models.Data{
			CodeError:     401,
			MessagesError: "Méthode non autorisée",
		}
		Error(w, r, &data)
	}
}

func RegisterHandle(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		render(w, r, "register.html", nil)
	case http.MethodPost:
		s := ""
		user := models.User{
			Username: r.FormValue("username"),
			Email:    r.FormValue("email"),
			Bio:      &s,
		}

		password := r.FormValue("password")
		vPassword := r.FormValue("confirm_password")

		if password == "" {
			data := models.Data{
				MessagesError: "Mot de passe requis",
			}
			render(w, r, "register.html", data)
			return
		}

		if password != vPassword {
			data := models.Data{
				MessagesError: "Les mots de passe ne correspondent pas",
			}
			render(w, r, "register.html", data)
			return
		}

		hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
		if err != nil {
			data := models.Data{
				CodeError:     500,
				MessagesError: "Erreur lors du chiffrement du mot de passe",
			}
			Error(w, r, &data)
			return
		}

		user.PasswordHash = string(hash)

		if err := repositoriespkg.CreateUser(&user); err != nil {
			if errors.Is(err, repositoriespkg.ErrEmailAlreadyExists) {
				data := models.Data{
					MessagesError: "email déjà utilisé",
				}
				render(w, r, "register.html", data)
				return
			}
			data := models.Data{
				CodeError:     400,
				MessagesError: "Erreur lors de la création du compte",
			}
			Error(w, r, &data)
			return
		}

		http.Redirect(w, r, "/auth/login", http.StatusSeeOther)

	default:
		data := models.Data{
			CodeError:     401,
			MessagesError: "Méthode non autorisée",
		}
		Error(w, r, &data)
	}
}

func LogoutHandle(w http.ResponseWriter, r *http.Request) {
	authpkg.DropSession(w, r)
	http.Redirect(w, r, "/", http.StatusSeeOther)
}

func Error(w http.ResponseWriter, r *http.Request, data *models.Data) {
	render(w, r, "error.html", data)
}
