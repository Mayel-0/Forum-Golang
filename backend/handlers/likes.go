package handlers

import (
	"lyrics/auth"
	"net/http"
)

/*var tpl *template.Template

func SetTemplates(t *template.Template) {
	tpl = t
} */

func LikeHandler(w http.ResponseWriter, r *http.Request) {
	userID, ok := auth.GetUserID(r)
	if !ok {
		http.Error(w, "unauthorized", http.StatusUnauthorized)
		return
	}
	_ = userID

	idLike := r.URL.Query().Get("id")
	if idLike == "" {
		http.Error(w, "id manquant", http.StatusBadGateway)
		return
	}

	switch r.Method {
	case http.MethodPost:

	default:
		http.Error(w, "méthode non autorisée", http.StatusMethodNotAllowed)
	}
}
