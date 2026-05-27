package handlers

import (
	"io"
	"lyrics/auth"
	dbpkg "lyrics/db"
	"lyrics/models"
	repositoriespkg "lyrics/repositories"
	"net/http"
	"path/filepath"
	"strconv"
	"strings"
	"text/template"
)

var templates map[string]*template.Template

func SetTemplates(t map[string]*template.Template) {
	templates = t
}

func render(w http.ResponseWriter, r *http.Request, name string, data any) {
	t, ok := templates[name]
	if !ok {
		http.Error(w, "Template introuvable: "+name, http.StatusInternalServerError)
		return
	}
	if err := t.ExecuteTemplate(w, "template.html", data); err != nil {
		http.Error(w, "Erreur lors du rendu du template", http.StatusInternalServerError)
		return
	}
}

func ForumIndexHandle(w http.ResponseWriter, r *http.Request) {
	http.Redirect(w, r, "/category/rock", http.StatusSeeOther)
}

func PublicProfileHandle(w http.ResponseWriter, r *http.Request) {
	username := r.PathValue("username")
	profileUser, err := repositoriespkg.GetUserByUsername(username)
	if err != nil {
		data := models.Data{
			CodeError:     404,
			MessagesError: "Utilisateur introuvable",
		}
		Error(w, r, &data)
		return
	}
	data := models.Data{ProfileUser: profileUser}
	if current, ok := auth.GetCurrentUser(r); ok {
		data.User = current
	}
	render(w, r, "user.html", data)
}

func ProfileHandle(w http.ResponseWriter, r *http.Request) {
	UserID, ok := auth.GetUserID(r)
	if !ok {
		data := models.Data{
			CodeError:     401,
			MessagesError: "Utilisateur non authentifié",
		}
		Error(w, r, &data)
		return
	}
	User, err := auth.GetUserByID(UserID)
	if err != nil {
		data := models.Data{
			CodeError:     500,
			MessagesError: "Erreur lors de la récupération de l'utilisateur",
		}
		Error(w, r, &data)
		return
	}

	switch r.Method {
	case http.MethodGet:
		render(w, r, "profile.html", models.Data{User: User})
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
			data := models.Data{
				CodeError:     500,
				MessagesError: "Erreur lors de la modification de l'utilisateur",
			}
			Error(w, r, &data)
			return
		}
		http.Redirect(w, r, "/profile", http.StatusSeeOther)
	default:
		data := models.Data{
			CodeError:     401,
			MessagesError: "Méthode non autorisée",
		}
		Error(w, r, &data)
	}
}

func UploadAvatarHandler(w http.ResponseWriter, r *http.Request) {
	UserID, ok := auth.GetUserID(r)
	if !ok {
		data := models.Data{
			CodeError:     401,
			MessagesError: "Utilisateur non authentifié",
		}
		Error(w, r, &data)
		return
	}

	// Lire le fichier uploadé
	r.ParseMultipartForm(5 << 20) // 5MB max
	file, header, err := r.FormFile("avatar")
	if err != nil {
		data := models.Data{
			CodeError:     400,
			MessagesError: "Fichier invalide",
		}
		Error(w, r, &data)
		return
	}
	defer file.Close()

	// Vérifier le type MIME
	ext := filepath.Ext(header.Filename)
	allowed := map[string]bool{".jpg": true, ".jpeg": true, ".png": true, ".webp": true}
	if !allowed[ext] {
		data := models.Data{
			CodeError:     400,
			MessagesError: "Format non supporté",
		}
		Error(w, r, &data)
		return
	}

	// Lire le contenu
	fileBytes, err := io.ReadAll(file)
	if err != nil {
		data := models.Data{
			CodeError:     500,
			MessagesError: "Erreur lecture fichier",
		}
		Error(w, r, &data)
		return
	}

	// Upload vers Supabase Storage
	avatarURL, err := repositoriespkg.UploadAvatar(UserID, ext, fileBytes)
	if err != nil {
		data := models.Data{
			CodeError:     500,
			MessagesError: "Erreur upload",
		}
		Error(w, r, &data)
		return
	}

	// Mettre à jour avatar_url en DB
	if err := repositoriespkg.UpdateUserAvatar(UserID, avatarURL); err != nil {
		data := models.Data{
			CodeError:     500,
			MessagesError: "Erreur mise à jour",
		}
		Error(w, r, &data)
		return
	}

	http.Redirect(w, r, "/profile", http.StatusSeeOther)
}

func UploadBanniereHandler(w http.ResponseWriter, r *http.Request) {
	UserID, ok := auth.GetUserID(r)
	if !ok {
		data := models.Data{
			CodeError:     401,
			MessagesError: "Utilisateur non authentifié",
		}
		Error(w, r, &data)
		return
	}

	r.ParseMultipartForm(5 << 20) // 5MB max
	file, header, err := r.FormFile("banniere")
	if err != nil {
		data := models.Data{
			CodeError:     400,
			MessagesError: "Fichier invalide",
		}
		Error(w, r, &data)
		return
	}
	defer file.Close()

	// Vérifier le type MIME
	ext := filepath.Ext(header.Filename)
	allowed := map[string]bool{".jpg": true, ".jpeg": true, ".png": true, ".webp": true}
	if !allowed[ext] {
		data := models.Data{
			CodeError:     400,
			MessagesError: "Format non supporté",
		}
		Error(w, r, &data)
		return
	}

	// Lire le contenu
	fileBytes, err := io.ReadAll(file)
	if err != nil {
		data := models.Data{
			CodeError:     500,
			MessagesError: "Erreur lecture fichier",
		}
		Error(w, r, &data)
		return
	}

	// Upload vers Supabase Storage
	avatarURL, err := repositoriespkg.UploadBanniere(UserID, ext, fileBytes)
	if err != nil {
		data := models.Data{
			CodeError:     500,
			MessagesError: "Erreur upload",
		}
		Error(w, r, &data)
		return
	}

	// Mettre à jour avatar_url en DB
	if err := repositoriespkg.UpdateUserBanniere(UserID, avatarURL); err != nil {
		data := models.Data{
			CodeError:     500,
			MessagesError: "Erreur mise à jour",
		}
		Error(w, r, &data)
		return
	}

	http.Redirect(w, r, "/profile", http.StatusSeeOther)
}

type SearchHandler struct {
	Templates map[string]*template.Template
}

func (h *SearchHandler) Search(w http.ResponseWriter, r *http.Request) {
	query := strings.TrimSpace(r.URL.Query().Get("q"))

	// Requête vide → redirect accueil
	if query == "" {
		http.Redirect(w, r, "/", http.StatusFound)
		return
	}

	// Pagination
	page := 1
	if p := r.URL.Query().Get("page"); p != "" {
		if parsed, err := strconv.Atoi(p); err == nil && parsed > 0 {
			page = parsed
		}
	}

	const perPage = 20

	posts, _, err := repositoriespkg.SearchPosts(dbpkg.Db, query, page, perPage)
	if err != nil {
		data := models.Data{
			CodeError:     500,
			MessagesError: "Erreur serveur",
		}
		Error(w, r, &data)
		return
	}

	data := models.Data{
		Query: query,
		Posts: posts,
	}

	tmpl, ok := h.Templates["search.html"]
	if !ok {
		data := models.Data{
			CodeError:     500,
			MessagesError: "Template introuvable",
		}
		Error(w, r, &data)
		return
	}

	if err := tmpl.ExecuteTemplate(w, "template.html", data); err != nil {
	}
}
