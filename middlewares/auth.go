package middlewares

import (
	"net/http"

	"github.com/gorilla/sessions"
)

type contextKey string

const UserKey contextKey = "username"

type authMiddleware struct {
	session sessions.Store
}

func NewAuthMiddleware(session sessions.Store) *authMiddleware {
	return &authMiddleware{
		session: session,
	}
}

func (m *authMiddleware) Authentication(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		session, _ := m.session.Get(r, "connections")
		_, ok := session.Values["username"].(string)

		if !ok {
			http.Redirect(w, r, "/login", http.StatusUnauthorized)
			return
		}

		next.ServeHTTP(w, r)
	})
}
