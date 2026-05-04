package handlers

import (
	"log"
	"lyrics/auth"
	"lyrics/db"
	"lyrics/models"
	repo "lyrics/repositories"
	"net/http"
	"regexp"
	"strings"

	"github.com/google/uuid"
)

func slugify(s string) string {
	// Convertir en minuscules
	s = strings.ToLower(s)
	// Remplacer les espaces par des tirets
	s = strings.ReplaceAll(s, " ", "-")
	// Supprimer les caractères spéciaux, garder seulement lettres, chiffres et tirets
	reg := regexp.MustCompile("[^a-z0-9-]+")
	s = reg.ReplaceAllString(s, "")
	// Supprimer les tirets multiples
	reg2 := regexp.MustCompile("-+")
	s = reg2.ReplaceAllString(s, "-")
	// Supprimer les tirets au début et à la fin
	s = strings.Trim(s, "-")
	return s
}

func PosteCreateHandler(w http.ResponseWriter, r *http.Request) {
	UserID, ok := auth.GetUserID(r)
	if !ok {
		http.Error(w, "Utilisateur non authentifié", http.StatusUnauthorized)
		return
	}
	switch r.Method {
	case http.MethodGet:
		categories, err := repo.GetAllCategories()
		if err != nil {
			http.Error(w, "Erreur lors de la récupération des catégories", http.StatusInternalServerError)
			log.Printf("Erreur récupération catégories: %v", err)
			return
		}
		render(w, "post-create.html", categories)
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
			Slug:     r.FormValue("category") + r.FormValue("subcategory") + "/" + slugify(r.FormValue("title")),
		}

		if err := db.Db.Create(&Post).Error; err != nil {
			http.Error(w, "Erreur lors de la création du post", http.StatusInternalServerError)
			log.Printf("Erreur création post: %v", err)
			return
		}

		http.Redirect(w, r, Post.Slug, http.StatusSeeOther)
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
	case http.MethodGet:
		render(w, "post-delete.html", nil)
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
	case http.MethodGet:
		render(w, "post-edit.html", nil)
	case http.MethodPost:
		userIDParsed, err := uuid.Parse(UserID)
		if err != nil {
			http.Error(w, "ID utilisateur invalide", http.StatusBadRequest)
			return
		}

		Post := models.Post{
			Title:    r.FormValue("title"),
			Body:     r.FormValue("body"),
			AuthorID: userIDParsed,
		}

		lastSulg := r.FormValue("slug")

		Newslug, err := repo.UpdatePostSlugFromTitle(lastSulg, Post.Title)
		if err != nil {
			http.Error(w, "Erreur lors de la modification", http.StatusInternalServerError)
			log.Printf("Erreur: %v", err)
			return
		}

		Post.Slug = Newslug

		if err := repo.UpdatePoste(&Post); err != nil {
			http.Error(w, "Erreur lors de la modification du post", http.StatusInternalServerError)
			log.Printf("Erreur modification post: %v", err)
			return
		}

		http.Redirect(w, r, Post.Slug, http.StatusSeeOther)
	default:
		http.Error(w, "Méthode non autorisée", http.StatusMethodNotAllowed)
	}
}
