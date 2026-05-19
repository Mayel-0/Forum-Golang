package handlers

import (
	"encoding/json"
	"lyrics/auth"
	"lyrics/models"
	"net/http"

	repo "lyrics/repositories"

	"github.com/google/uuid"
)

func AddCommentHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodPost:
		UserId, ok := auth.GetUserID(r)
		if !ok {
			data := models.Data{
				CodeError:     401,
				MessagesError: "Utilisateur non authentifié",
			}
			Error(w, r, &data)
			return
		}
		UserIdParsed, err := uuid.Parse(UserId)
		if err != nil {
			data := models.Data{
				CodeError:     401,
				MessagesError: "ID utilisateur invalide",
			}
			Error(w, r, &data)
			return
		}

		postID, err := uuid.Parse(r.FormValue("post_id"))
		if err != nil {
			data := models.Data{
				CodeError:     401,
				MessagesError: "ID post invalide",
			}
			Error(w, r, &data)
			return
		}

		var parentID *uuid.UUID
		if raw := r.FormValue("parent_id"); raw != "" {
			if pid, err := uuid.Parse(raw); err == nil {
				parentID = &pid
			}
		}

		comment := models.Comment{
			AuthorID: UserIdParsed,
			PostID:   postID,
			Body:     r.FormValue("content"),
			ParentID: parentID,
		}
		if err := repo.AddComment(&comment); err != nil {
			data := models.Data{
				CodeError:     500,
				MessagesError: "Erreur lors de la création du commentaire",
			}
			Error(w, r, &data)
			return
		}

		if r.Header.Get("X-Requested-With") == "fetch" {
			user, _ := auth.GetCurrentUser(r)
			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(map[string]any{
				"id":                comment.ID,
				"body":              comment.Body,
				"created_at":        comment.CreatedAt,
				"author_username":   user.Username,
				"author_avatar_url": user.AvatarURL,
			})
			return
		}

		http.Redirect(w, r, r.FormValue("slug"), http.StatusSeeOther)
	default:
		data := models.Data{
			CodeError:     401,
			MessagesError: "Méthode non autorisée",
		}
		Error(w, r, &data)
	}
}

func DeleteCommentHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodPost:
		UserId, ok := auth.GetUserID(r)
		if !ok {
			data := models.Data{
				CodeError:     401,
				MessagesError: "Utilisateur non authentifié",
			}
			Error(w, r, &data)
			return
		}

		UserIdParsed, err := uuid.Parse(UserId)
		if err != nil {
			data := models.Data{
				CodeError:     401,
				MessagesError: "ID utilisateur invalide",
			}
			Error(w, r, &data)
			return
		}

		CommentId, err := uuid.Parse(r.FormValue("comment_id"))
		if err != nil {
			data := models.Data{
				CodeError:     401,
				MessagesError: "ID commentaire invalide",
			}
			Error(w, r, &data)
			return
		}

		comment, err := repo.GetCommentByID(CommentId)
		if err != nil {
			data := models.Data{
				CodeError:     404,
				MessagesError: "Commentaire non trouvé",
			}
			Error(w, r, &data)
			return
		}

		if !auth.VerifyUserRequest(UserIdParsed, comment.AuthorID) {
			data := models.Data{
				CodeError:     403,
				MessagesError: "Utilisateur non autorisé à modifier ce commentaire",
			}
			Error(w, r, &data)
			return
		}

		if err := repo.DeleteComment(CommentId); err != nil {
			data := models.Data{
				CodeError:     500,
				MessagesError: "Erreur lors de la suppression du commentaire",
			}
			Error(w, r, &data)
			return
		}

		http.Redirect(w, r, r.FormValue("slug"), http.StatusSeeOther)

	default:
		data := models.Data{
			CodeError:     401,
			MessagesError: "Méthode non autorisée",
		}
		Error(w, r, &data)
	}
}
