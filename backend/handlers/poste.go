package handlers

import (
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
		data := models.Data{
			CodeError:     401,
			MessagesError: "Utilisateur non authentifié",
		}
		Error(w, r, &data)
		return
	}
	switch r.Method {
	case http.MethodGet:
		categories, err := repo.GetAllCategories()
		if err != nil {
			data := models.Data{
				CodeError:     500,
				MessagesError: "Erreur lors de la récupération des catégories",
			}
			Error(w, r, &data)
			return
		}

		data := models.Data{Categories: categories}
		if user, ok := auth.GetCurrentUser(r); ok {
			data.User = user
		}
		render(w,r, "post-create.html", data)
	case http.MethodPost:
		userIDParsed, err := uuid.Parse(UserID)
		if err != nil {
			data := models.Data{
				CodeError:     401,
				MessagesError: "ID utilisateur invalide",
			}
			Error(w, r, &data)
			return
		}

		categorySlug := sanitizeSlugPart(r.FormValue("category"))
		subcategoryRaw := sanitizeSlugPart(r.FormValue("subcategory"))
		if categorySlug == "" || subcategoryRaw == "" {
			data := models.Data{
				CodeError:     401,
				MessagesError: "Catégorie ou sous-catégorie invalide",
			}
			Error(w, r, &data)
			return
		}
		if !validSubCategory(subcategoryRaw) {
			data := models.Data{
				CodeError:     401,
				MessagesError: "Catégorie ou sous-catégorie invalide",
			}
			Error(w, r, &data)
			return
		}

		categoryID, err := repo.GetCategoryIDBySlug(r.FormValue("category"))
		if err != nil {
			data := models.Data{
				CodeError:     401,
				MessagesError: "Catégorie ou sous-catégorie invalide",
			}
			Error(w, r, &data)
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
			data := models.Data{
				CodeError:     500,
				MessagesError: "Erreur lors de la création du post",
			}
			Error(w, r, &data)
			return
		}

		reelSlug := "/p/" + Post.Slug

		http.Redirect(w, r, reelSlug, http.StatusSeeOther)
	default:
		data := models.Data{
			CodeError:     401,
			MessagesError: "Méthode non autorisée",
		}
		Error(w, r, &data)
	}
}

func PosteDeleteHandler(w http.ResponseWriter, r *http.Request) {
	UserID, ok := auth.GetUserID(r)
	if !ok {
		data := models.Data{
			CodeError:     401,
			MessagesError: "Utilisateur non authentifié",
		}
		Error(w, r, &data)
		return
	}
	switch r.Method {
	case http.MethodPost:
		userIDParsed, err := uuid.Parse(UserID)
		if err != nil {
			data := models.Data{
				CodeError:     401,
				MessagesError: "ID utilisateur invalide",
			}
			Error(w, r, &data)
			return
		}

		postID := r.FormValue("postID")
		postIDParsed, err := uuid.Parse(postID)
		if err != nil {
			data := models.Data{
				CodeError:     401,
				MessagesError: "ID post invalide",
			}
			Error(w, r, &data)
			return
		}

		post, err := repo.GetPostByID(postIDParsed)
		if err != nil {
			data := models.Data{
				CodeError:     404,
				MessagesError: "Post introuvable",
			}
			Error(w, r, &data)
			return
		}

		if !auth.VerifyUserRequest(userIDParsed, post.AuthorID) {
			data := models.Data{
				CodeError:     401,
				MessagesError: "Utilisateur non autorisé à modifier ce post",
			}
			Error(w, r, &data)
			return
		}

		if err := repo.DeletePoste(&post); err != nil {
			data := models.Data{
				CodeError:     500,
				MessagesError: "Erreur lors de la suppression du post",
			}
			Error(w, r, &data)
			return
		}

		http.Redirect(w, r, "acceuil.html", http.StatusSeeOther)
	default:
		data := models.Data{
			CodeError:     401,
			MessagesError: "Méthode non autorisée",
		}
		Error(w, r, &data)
	}
}

func PosteModifierHandle(w http.ResponseWriter, r *http.Request) {
	UserID, ok := auth.GetUserID(r)
	if !ok {
		data := models.Data{
			CodeError:     401,
			MessagesError: "Utilisateur non authentifié",
		}
		Error(w, r, &data)
		return
	}
	switch r.Method {
	case http.MethodGet:
		render(w, r, "post-edit.html", nil)
	case http.MethodPost:
		userIDParsed, err := uuid.Parse(UserID)
		if err != nil {
			data := models.Data{
				CodeError:     401,
				MessagesError: "ID utilisateur invalide",
			}
			Error(w, r, &data)
			return
		}

		postID := r.FormValue("postID")
		postIDParsed, err := uuid.Parse(postID)
		if err != nil {
			data := models.Data{
				CodeError:     401,
				MessagesError: "ID post invalide",
			}
			Error(w, r, &data)
			return
		}

		post, err := repo.GetPostByID(postIDParsed)
		if err != nil {
			data := models.Data{
				CodeError:     404,
				MessagesError: "Post introuvable",
			}
			Error(w, r, &data)
			return
		}

		if !auth.VerifyUserRequest(userIDParsed, post.AuthorID) {
			data := models.Data{
				CodeError:     401,
				MessagesError: "Utilisateur non autorisé à modifier ce post",
			}
			Error(w, r, &data)
			return
		}

		post.Title = r.FormValue("title")
		post.Body = r.FormValue("body")

		newSlug, err := repo.UpdatePostSlugFromTitle(post.Slug, post.Title)
		if err != nil {
			data := models.Data{
				CodeError:     500,
				MessagesError: "Erreur lors de la modification",
			}
			Error(w, r, &data)
			return
		}
		post.Slug = newSlug

		if err := repo.UpdatePoste(&post); err != nil {
			data := models.Data{
				CodeError:     500,
				MessagesError: "Erreur lors de la modification du post",
			}
			Error(w, r, &data)
			return
		}

		http.Redirect(w, r, post.Slug, http.StatusSeeOther)
	default:
		data := models.Data{
			CodeError:     401,
			MessagesError: "Méthode non autorisée",
		}
		Error(w, r, &data)
	}
}

func loadAsideData(data *models.Data) {
	if top, err := repositoriespkg.GetTopUsers(3); err == nil {
		data.TopUsers = top
	}
	if recent, err := repo.GetRecentPosts(5); err == nil {
		data.RecentPosts = recent
	}
	if popular, err := repo.GetPopularPosts(5); err == nil {
		data.PopularPosts = popular
	}
}

func PostShowHandle(w http.ResponseWriter, r *http.Request) {
	slug := r.PathValue("slug")

	post, err := repo.GetPostBySlug(slug)
	if err != nil {
		data := models.Data{
			CodeError:     404,
			MessagesError: "Post non trouvé",
		}
		Error(w, r, &data)
		return
	}

	comments, err := repo.GetCommentsByPostID(post.ID)
	if err != nil {
		data := models.Data{
			CodeError:     404,
			MessagesError: "Commentaire non trouvé",
		}
		Error(w, r, &data)
		return
	}

	data := models.Data{Post: post, Comments: comments}
	loadAsideData(&data)

	if user, ok := auth.GetCurrentUser(r); ok {
		data.User = user
		data.UserLikedPost = repo.HasUserLikedPost(user.ID, post.ID)
		data.UserLikedCommentIDs = repo.GetUserLikedCommentIDs(user.ID, post.ID)
	}

	render(w,r, "post.html", data)
}

func SubCategoryHandle(w http.ResponseWriter, r *http.Request) {
	slug := "/" + r.PathValue("slug")
	subRaw := r.PathValue("subcategory")

	if !validSubCategory(subRaw) {
		data := models.Data{
			CodeError:     404,
			MessagesError: "Sous categorie non trouvé",
		}
		Error(w, r, &data)
		return
	}
	sub := models.PostSubCategory(subRaw)

	categoryID, err := repo.GetCategoryIDBySlug(slug)
	if err != nil {
		data := models.Data{
			CodeError:     404,
			MessagesError: "Sous categorie non trouvé",
		}
		Error(w, r, &data)
		return
	}

	posts, err := repo.GetAllPostsByCategoryAndSub(categoryID, sub)
	if err != nil {
		data := models.Data{
			CodeError:     500,
			MessagesError: "Erreur récupération des posts",
		}
		Error(w, r, &data)
		return
	}

	category, err := repo.GetCategoryByID(categoryID)
	if err != nil {
		data := models.Data{
			CodeError:     500,
			MessagesError: "Erreur récupération de la catégorie",
		}
		Error(w, r, &data)
		return
	}

	topUsers, err := repositoriespkg.GetTopUsers(3)
	if err != nil {
		data := models.Data{
			CodeError:     500,
			MessagesError: "Erreur récupération des utilisateurs",
		}
		Error(w, r, &data)
		return
	}

	labels := map[string]string{
		"artistes":   "Artistes",
		"concerts":   "Concerts",
		"nouveautes": "Nouveautés",
	}

	data := models.Data{
		Posts:            posts,
		Category:         category,
		TopUsers:         topUsers,
		SubCategoryLabel: labels[subRaw],
	}
	loadAsideData(&data)
	if user, ok := auth.GetCurrentUser(r); ok {
		data.User = user
	}

	render(w,r, "subcategory.html", data)
}

func CategoryHandle(w http.ResponseWriter, r *http.Request) {
	slug := "/" + r.PathValue("slug")

	categoryID, err := repo.GetCategoryIDBySlug(slug)
	if err != nil {
		data := models.Data{
			CodeError:     404,
			MessagesError: "Catégorie non trouvée",
		}
		Error(w, r, &data)
		return
	}

	postsA, err := repo.GetAllPostsByCategoryLimitArtistes(categoryID)
	if err != nil {
		data := models.Data{
			CodeError:     500,
			MessagesError: "Erreur récupération des posts",
		}
		Error(w, r, &data)
		return
	}

	postsC, err := repo.GetAllPostsByCategoryLimitConcerts(categoryID)
	if err != nil {
		data := models.Data{
			CodeError:     500,
			MessagesError: "Erreur récupération des posts",
		}
		Error(w, r, &data)
		return
	}

	postsN, err := repo.GetAllPostsByCategoryLimitNouveautes(categoryID)
	if err != nil {
		data := models.Data{
			CodeError:     500,
			MessagesError: "Erreur récupération des posts",
		}
		Error(w, r, &data)
		return
	}

	topUsers, err := repositoriespkg.GetTopUsers(3)
	if err != nil {
		data := models.Data{
			CodeError:     500,
			MessagesError: "Erreur récupération des utilisateurs",
		}
		Error(w, r, &data)
		return
	}

	category, err := repo.GetCategoryByID(categoryID)
	if err != nil {
		data := models.Data{
			CodeError:     500,
			MessagesError: "Erreur récupération de la catégorie",
		}
		Error(w, r, &data)
		return
	}

	totalA, _ := repo.CountPostsByCategoryAndSub(categoryID, models.SubCategoryArtistes)
	totalC, _ := repo.CountPostsByCategoryAndSub(categoryID, models.SubCategoryConcerts)
	totalN, _ := repo.CountPostsByCategoryAndSub(categoryID, models.SubCategoryNouveautes)

	data := models.Data{
		TopUsers:        topUsers,
		PostsArtists:    postsA,
		PostsConcerts:   postsC,
		PostsNouveautes: postsN,
		Category:        category,
		TotalArtistes:   totalA,
		TotalConcerts:   totalC,
		TotalNouveautes: totalN,
	}
	loadAsideData(&data)
	if user, ok := auth.GetCurrentUser(r); ok {
		data.User = user
	}

	render(w,r, "index.html", data)
}
