package internal

import (
	"net/http"

	"github.com/gorilla/sessions"
)

func ClearSession(session sessions.Store, w http.ResponseWriter, r *http.Request) {
	currentSession, _ := session.Get(r, "connections")
	currentSession.Options.MaxAge = -1
	currentSession.Save(r, w)
}

func CreateSession(session sessions.Store, w http.ResponseWriter, r *http.Request, username string) {
	currentSession, _ := session.Get(r, "connections")
	currentSession.Values["username"] = username
	currentSession.Save(r, w)
}
