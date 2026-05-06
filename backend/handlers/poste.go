package handlers

import (
	"log"
	auth "lyrics/auth"
	"lyrics/db"
	"lyrics/models"
	repo "lyrics/repositories"
	repositoriespkg "lyrics/repositories"
	"net/http"
	"regexp"
	"strings"
	"text/template"

	"github.com/google/uuid"
)

var tpl map[string]*template.Template

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

func sanitizeSlugPart(value string) string {
	value = strings.TrimSpace(value)
	value = strings.TrimPrefix(value, "/")
	value = strings.TrimSuffix(value, "/")
	return strings.TrimSpace(value)
}

func validSubCategory(value string) bool {
	switch value {
	case string(models.SubCategoryConcerts), string(models.SubCategoryArtistes), string(models.SubCategoryNouveautes):
		return true
	}
	return false
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

		data := models.Data{Categories: categories}
		if user, ok := auth.GetCurrentUser(r); ok {
			data.User = user
		}
		render(w, "post-create.html", data)
	case http.MethodPost:
		userIDParsed, err := uuid.Parse(UserID)
		if err != nil {
			http.Error(w, "ID utilisateur invalide", http.StatusBadRequest)
			return
		}

		categorySlug := sanitizeSlugPart(r.FormValue("category"))
		subcategoryRaw := sanitizeSlugPart(r.FormValue("subcategory"))
		if categorySlug == "" || subcategoryRaw == "" {
			http.Error(w, "Catégorie ou sous-catégorie invalide", http.StatusBadRequest)
			return
		}
		if !validSubCategory(subcategoryRaw) {
			http.Error(w, "Sous-catégorie invalide", http.StatusBadRequest)
			return
		}

		categoryID, err := repo.GetCategoryIDBySlug(r.FormValue("category"))
		if err != nil {
			http.Error(w, "Catégorie invalide", http.StatusBadRequest)
			return
		}

		Post := models.Post{
			AuthorID:    userIDParsed,
			CategoryID:  &categoryID,
			Title:       r.FormValue("title"),
			Body:        r.FormValue("body"),
			Slug:        categorySlug + "/" + subcategoryRaw + "/" + slugify(r.FormValue("title")),
			SubCategory: models.PostSubCategory(subcategoryRaw),
		}

		if err := db.Db.Create(&Post).Error; err != nil {
			http.Error(w, "Erreur lors de la création du post", http.StatusInternalServerError)
			log.Printf("Erreur création post: %v", err)
			return
		}

		reelSlug := "/p/" + Post.Slug

		http.Redirect(w, r, reelSlug, http.StatusSeeOther)
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
		userIDParsed, err := uuid.Parse(UserID)
		if err != nil {
			http.Error(w, "ID utilisateur invalide", http.StatusBadRequest)
			return
		}

		postID := r.FormValue("postID")
		postIDParsed, err := uuid.Parse(postID)
		if err != nil {
			http.Error(w, "ID post invalide", http.StatusBadRequest)
			return
		}

		post, err := repo.GetPostByID(postIDParsed)
		if err != nil {
			http.Error(w, "Post introuvable", http.StatusNotFound)
			log.Printf("Erreur récupération post: %v", err)
			return
		}

		if !auth.VerifyUserRequest(userIDParsed, post.AuthorID) {
			http.Error(w, "Utilisateur non autorisé à modifier ce post", http.StatusForbidden)
			return
		}

		if err := repo.DeletePoste(&post); err != nil {
			http.Error(w, "Erreur lors de la suppression du post", http.StatusInternalServerError)
			log.Printf("Erreur suppression post: %v", err)
			return
		}

		http.Redirect(w, r, "acceuil.html", http.StatusSeeOther)
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

		postID := r.FormValue("postID")
		postIDParsed, err := uuid.Parse(postID)
		if err != nil {
			http.Error(w, "ID post invalide", http.StatusBadRequest)
			return
		}

		post, err := repo.GetPostByID(postIDParsed)
		if err != nil {
			http.Error(w, "Post introuvable", http.StatusNotFound)
			log.Printf("Erreur récupération post: %v", err)
			return
		}

		if !auth.VerifyUserRequest(userIDParsed, post.AuthorID) {
			http.Error(w, "Utilisateur non autorisé à modifier ce post", http.StatusForbidden)
			return
		}

		post.Title = r.FormValue("title")
		post.Body = r.FormValue("body")

		newSlug, err := repo.UpdatePostSlugFromTitle(post.Slug, post.Title)
		if err != nil {
			http.Error(w, "Erreur lors de la modification", http.StatusInternalServerError)
			log.Printf("Erreur: %v", err)
			return
		}
		post.Slug = newSlug

		if err := repo.UpdatePoste(&post); err != nil {
			http.Error(w, "Erreur lors de la modification du post", http.StatusInternalServerError)
			log.Printf("Erreur modification post: %v", err)
			return
		}

		http.Redirect(w, r, post.Slug, http.StatusSeeOther)
	default:
		http.Error(w, "Méthode non autorisée", http.StatusMethodNotAllowed)
	}
}

func PostShowHandle(w http.ResponseWriter, r *http.Request) {
	slug := r.PathValue("slug")

	post, err := repo.GetPostBySlug(slug)
	if err != nil {
		http.Error(w, "Post non trouvé", http.StatusNotFound)
		return
	}

	data := models.Data{Post: post}

	if user, ok := auth.GetCurrentUser(r); ok {
		data.User = user
	}

	render(w, "post.html", data)
}

func CategoryHandle(w http.ResponseWriter, r *http.Request) {
	slug := "/" + r.PathValue("slug")

	categoryID, err := repo.GetCategoryIDBySlug(slug)
	if err != nil {
		http.Error(w, "Catégorie non trouvée", http.StatusNotFound)
		return
	}

	postsA, err := repo.GetAllPostsByCategoryLimitArtistes(categoryID)
	if err != nil {
		http.Error(w, "Erreur récupération des posts", http.StatusInternalServerError)
		return
	}

	postsC, err := repo.GetAllPostsByCategoryLimitConcerts(categoryID)
	if err != nil {
		http.Error(w, "Erreur récupération des posts", http.StatusInternalServerError)
		return
	}

	postsN, err := repo.GetAllPostsByCategoryLimitNouveautes(categoryID)
	if err != nil {
		http.Error(w, "Erreur récupération des posts", http.StatusInternalServerError)
		return
	}

	topUsers, err := repositoriespkg.GetTopUsers(3)
	if err != nil {
		http.Error(w, "Erreur récupération des utilisateurs", http.StatusInternalServerError)
		return
	}

	category, err := repo.GetCategoryByID(categoryID)
	if err != nil {
		http.Error(w, "Erreur récupération de la catégorie", http.StatusInternalServerError)
		return
	}

	data := models.Data{
		TopUsers:        topUsers,
		PostsArtists:    postsA,
		PostsConcerts:   postsC,
		PostsNouveautes: postsN,
		Category:        category,
	}
	if user, ok := auth.GetCurrentUser(r); ok {
		data.User = user
	}

	render(w, "index.html", data)
}
