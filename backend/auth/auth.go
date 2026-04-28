package auth

import (
	"context"
	"errors"
	"log"
	"lyrics/models"
	"net/http"
	"strings"

	dbpkg "lyrics/db"

	"github.com/gorilla/sessions"
)

type contextKey string

const (
	sessionName                 = "auth-session"
	userIDContextKey contextKey = "auth.user_id"
)

var Store *sessions.CookieStore

func SetStore(s *sessions.CookieStore) {
	Store = s
}

func SetSession(w http.ResponseWriter, r *http.Request, userID string) error {
	if Store == nil {
		return errors.New("session store not initialized")
	}

	session, err := Store.Get(r, sessionName)
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

func DropSession(w http.ResponseWriter, r *http.Request) error {
	if Store == nil {
		return errors.New("session store not initialized")
	}

	session, err := Store.Get(r, sessionName)
	if err != nil {
		return err
	}

	session.Options.MaxAge = -1
	return session.Save(r, w)
}

func ClearSession(w http.ResponseWriter, r *http.Request) error {
	if Store == nil {
		return errors.New("session store not initialized")
	}

	session, err := Store.Get(r, sessionName)
	if err != nil {
		return err
	}

	session.Options.MaxAge = -1
	return session.Save(r, w)
}

func GetUserID(r *http.Request) (string, bool) {
	userID, ok := r.Context().Value(userIDContextKey).(string)
	return userID, ok
}

func GetUserByID(userID string) (models.User, error) {
	var user models.User
	err := dbpkg.Db.Where("id = ?", userID).First(&user).Error
	if err != nil {
		return models.User{}, err
	}
	return user, nil
}

func RequireAuth(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if Store == nil {
			log.Print("session store not initialized")
			rejectUnauthorized(w, r)
			return
		}

		session, err := Store.Get(r, sessionName)
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
