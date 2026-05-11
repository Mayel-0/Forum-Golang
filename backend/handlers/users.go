package handlers

import (
	"io"
	"log"
	"lyrics/auth"
	"lyrics/models"
	repositoriespkg "lyrics/repositories"
	"net/http"
	"path/filepath"
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

	http.Redirect(w, r, "/category/rock", http.StatusSeeOther)
	// topUsers, err := repositoriespkg.GetTopUsers(3)
	//if err != nil {
	//	http.Error(w, "Erreur récupération des utilisateurs", http.StatusInternalServerError)
	//	return
	//}
	//data := models.Data{TopUsers: topUsers}
	//if user, ok := auth.GetCurrentUser(r); ok {
	//	data.User = user
	//}
	//render(w, "index.html", data)
}

func PublicProfileHandle(w http.ResponseWriter, r *http.Request) {
	username := r.PathValue("username")
	profileUser, err := repositoriespkg.GetUserByUsername(username)
	if err != nil {
		http.Error(w, "Utilisateur introuvable", http.StatusNotFound)
		return
	}
	data := models.Data{ProfileUser: profileUser}
	if current, ok := auth.GetCurrentUser(r); ok {
		data.User = current
	}
	render(w, "user.html", data)
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

func UploadAvatarHandler(w http.ResponseWriter, r *http.Request) {
	UserID, ok := auth.GetUserID(r)
	if !ok {
		http.Error(w, "Utilisateur non authentifié", http.StatusUnauthorized)
		return
	}

	// Lire le fichier uploadé
	r.ParseMultipartForm(5 << 20) // 5MB max
	file, header, err := r.FormFile("avatar")
	if err != nil {
		http.Error(w, "Fichier invalide", http.StatusBadRequest)
		return
	}
	defer file.Close()

	// Vérifier le type MIME
	ext := filepath.Ext(header.Filename)
	allowed := map[string]bool{".jpg": true, ".jpeg": true, ".png": true, ".webp": true}
	if !allowed[ext] {
		http.Error(w, "Format non supporté", http.StatusBadRequest)
		return
	}

	// Lire le contenu
	fileBytes, err := io.ReadAll(file)
	if err != nil {
		http.Error(w, "Erreur lecture fichier", http.StatusInternalServerError)
		return
	}

	// Upload vers Supabase Storage
	avatarURL, err := repositoriespkg.UploadAvatar(UserID, ext, fileBytes)
	if err != nil {
		http.Error(w, "Erreur upload", http.StatusInternalServerError)
		return
	}

	// Mettre à jour avatar_url en DB
	if err := repositoriespkg.UpdateUserAvatar(UserID, avatarURL); err != nil {
		http.Error(w, "Erreur mise à jour", http.StatusInternalServerError)
		return
	}

	http.Redirect(w, r, "/profile", http.StatusSeeOther)
}

func UploadBanniereHandler(w http.ResponseWriter, r *http.Request) {
	UserID, ok := auth.GetUserID(r)
	if !ok {
		http.Error(w, "Utilisateur non authentifié", http.StatusUnauthorized)
		return
	}

	r.ParseMultipartForm(5 << 20) // 5MB max
	file, header, err := r.FormFile("avatar")
	if err != nil {
		http.Error(w, "Fichier invalide", http.StatusBadRequest)
		return
	}
	defer file.Close()

	// Vérifier le type MIME
	ext := filepath.Ext(header.Filename)
	allowed := map[string]bool{".jpg": true, ".jpeg": true, ".png": true, ".webp": true}
	if !allowed[ext] {
		http.Error(w, "Format non supporté", http.StatusBadRequest)
		return
	}

	// Lire le contenu
	fileBytes, err := io.ReadAll(file)
	if err != nil {
		http.Error(w, "Erreur lecture fichier", http.StatusInternalServerError)
		return
	}

	// Upload vers Supabase Storage
	avatarURL, err := repositoriespkg.UploadBanniere(UserID, ext, fileBytes)
	if err != nil {
		http.Error(w, "Erreur upload", http.StatusInternalServerError)
		return
	}

	// Mettre à jour avatar_url en DB
	if err := repositoriespkg.UpdateUserBanniere(UserID, avatarURL); err != nil {
		http.Error(w, "Erreur mise à jour", http.StatusInternalServerError)
		return
	}

	http.Redirect(w, r, "/profile", http.StatusSeeOther)
}

func SearchHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:

	default:
		http.Error(w, "Méthode non autorisée", http.StatusMethodNotAllowed)
	}
}
