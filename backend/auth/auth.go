package auth

import (
	"context"
	"log"
	"net/http"
	"strings"

	"github.com/gorilla/sessions"
)

type contextKey string

const userIDContextKey contextKey = "auth.user_id"

var Store *sessions.CookieStore

func SetStore(s *sessions.CookieStore) {
	Store = s
}

func SetSession(w http.ResponseWriter, r *http.Request, userID string) error {
	session, err := Store.Get(r, "session-name")
	if err != nil {
		return err
	}

	session.Values["user_id"] = userID
	session.Options.MaxAge = 86400 * 7 // 7 jours
	session.Options.HttpOnly = true
	session.Options.Secure = r.TLS != nil || r.Header.Get("X-Forwarded-Proto") == "https"
	session.Options.SameSite = http.SameSiteLaxMode

	return session.Save(r, w)
}

func GetUserID(r *http.Request) (string, bool) {
	userID, ok := r.Context().Value(userIDContextKey).(string)
	return userID, ok
}

func RequireAuth(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		session, err := Store.Get(r, "session-name")
		if err != nil {
			log.Printf("Erreur session: %v", err)
			rejectUnauthorized(w, r)
			return
		}

		if session.IsNew {
			rejectUnauthorized(w, r)
			return
		}

		userID, ok := session.Values["user_id"].(string)
		if !ok || userID == "" {
			rejectUnauthorized(w, r)
			return
		}

		ctx := context.WithValue(r.Context(), userIDContextKey, userID)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func rejectUnauthorized(w http.ResponseWriter, r *http.Request) {
	log.Printf("requireAuth refusé: méthode=%s path=%s", r.Method, r.URL.Path)

	accept := r.Header.Get("Accept")
	if r.Method == http.MethodGet && strings.Contains(accept, "text/html") {
		http.Redirect(w, r, "/auth/login", http.StatusSeeOther)
		return
	}

	http.Error(w, "unauthorized", http.StatusUnauthorized)
}
