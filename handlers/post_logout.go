package handlers

import (
	"net/http"

	"github.com/gorilla/sessions"
	"github.com/ngsalvo/roadmapsh-personal-blog/internal"
	datastar "github.com/starfederation/datastar/sdk/go"
)

type PostLogout interface {
	Handle(w http.ResponseWriter, r *http.Request)
}

type postLogout struct {
	session sessions.Store
}

func NewPostLogout(session sessions.Store) PostLogout {
	return &postLogout{
		session: session,
	}
}

func (h *postLogout) Handle(w http.ResponseWriter, r *http.Request) {
	internal.ClearSession(h.session, w, r)

	sse := datastar.NewSSE(w, r)

	sse.Redirect("/")
}
