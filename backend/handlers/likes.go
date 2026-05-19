package handlers

import (
	"lyrics/auth"
	"lyrics/db"
	"lyrics/models"
	"net/http"

	"github.com/google/uuid"
)

/*var tpl *template.Template

func SetTemplates(t *template.Template) {
	tpl = t
} */

func LikeHandlerAdd(w http.ResponseWriter, r *http.Request) {
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
		PostIDStr := r.FormValue("post_id")
		CommentIDStr := r.FormValue("comment_id")

		userUUID, err := uuid.Parse(UserID)
		if err != nil {
			data := models.Data{
				CodeError:     401,
				MessagesError: "ID utilisateur invalide",
			}
			Error(w, r, &data)
			return
		}

		like := models.Likes{
			UserID: userUUID,
		}

		if PostIDStr != "" {
			PostUUID, err := uuid.Parse(PostIDStr)
			if err != nil {
				data := models.Data{
					CodeError:     401,
					MessagesError: "ID utilisateur invalide",
				}
				Error(w, r, &data)
				return
			}
			like.PostID = &PostUUID
		}

		if CommentIDStr != "" {
			CommentUUID, err := uuid.Parse(CommentIDStr)
			if err != nil {
				data := models.Data{
					CodeError:     401,
					MessagesError: "ID commentaire invalide",
				}
				Error(w, r, &data)
				return
			}
			like.CommentID = &CommentUUID
		}

		if err := db.Db.Create(&like).Error; err != nil {
			data := models.Data{
				CodeError:     500,
				MessagesError: "Erreur lors de la création du like",
			}
			Error(w, r, &data)
			return
		}

		w.WriteHeader(http.StatusCreated)

	default:
		data := models.Data{
			CodeError:     401,
			MessagesError: "Méthode non autorisée",
		}
		Error(w, r, &data)
	}
}

func LikeHandlerRm(w http.ResponseWriter, r *http.Request) {
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
		PostIDStr := r.FormValue("post_id")
		CommentIDStr := r.FormValue("comment_id")

		userUUID, err := uuid.Parse(UserID)
		if err != nil {
			data := models.Data{
				CodeError:     401,
				MessagesError: "ID utilisateur invalide",
			}
			Error(w, r, &data)
			return
		}

		query := db.Db.Where("user_id = ?", userUUID)

		if PostIDStr != "" {
			PostUUID, err := uuid.Parse(PostIDStr)
			if err != nil {
				data := models.Data{
					CodeError:     401,
					MessagesError: "ID utilisateur invalide",
				}
				Error(w, r, &data)
				return
			}
			query = query.Where("post_id = ?", PostUUID)
		}
		if CommentIDStr != "" {
			CommentUUID, err := uuid.Parse(CommentIDStr)
			if err != nil {
				data := models.Data{
					CodeError:     401,
					MessagesError: "ID commentaire invalide",
				}
				Error(w, r, &data)
				return
			}
			query = query.Where("comment_id = ?", CommentUUID)
		}

		if err := query.Delete(&models.Likes{}).Error; err != nil {
			data := models.Data{
				CodeError:     500,
				MessagesError: "Erreur lors de la suppression du like",
			}
			Error(w, r, &data)
			return
		}

		w.WriteHeader(http.StatusNoContent)

	default:
		data := models.Data{
			CodeError:     401,
			MessagesError: "Méthode non autorisée",
		}
		Error(w, r, &data)
	}
}
