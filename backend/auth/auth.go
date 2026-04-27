package auth

import (
	"context"
	"log"
	"net/http"
	"strings"
	"sync"
	"time"

	models "lyrics/models"
)

type contextKey string

const userIDContextKey contextKey = "auth.user_id"

var (
	sessionsMu sync.RWMutex
	sessions   = map[string]models.Session{}
)

func StoreSession(token string, session models.Session) {
	sessionsMu.Lock()
	defer sessionsMu.Unlock()
	sessions[token] = session
}

func DeleteSession(token string) {
	sessionsMu.Lock()
	defer sessionsMu.Unlock()
	delete(sessions, token)
}

func GetUserID(r *http.Request) (string, bool) {
	userID, ok := r.Context().Value(userIDContextKey).(string)
	return userID, ok
}

func RequireAuth(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		cookie, err := r.Cookie("session_token")
		if err != nil || cookie.Value == "" {
			rejectUnauthorized(w, r)
			return
		}

		sessionsMu.RLock()
		session, ok := sessions[cookie.Value]
		sessionsMu.RUnlock()
		if !ok {
			clearSessionCookie(w)
			rejectUnauthorized(w, r)
			return
		}

		if session.Userid == "" || time.Now().After(session.Expiry) {
			DeleteSession(cookie.Value)
			clearSessionCookie(w)
			rejectUnauthorized(w, r)
			return
		}

		ctx := context.WithValue(r.Context(), userIDContextKey, session.Userid)
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

func clearSessionCookie(w http.ResponseWriter) {
	http.SetCookie(w, &http.Cookie{
		Name:     "session_token",
		Value:    "",
		Path:     "/",
		MaxAge:   -1,
		HttpOnly: true,
		Secure:   true,
		SameSite: http.SameSiteLaxMode,
	})
}

func requireAuth(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		RequireAuth(next).ServeHTTP(w, r)
	}
}
