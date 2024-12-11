package handlers

import (
	"net/http"

	"github.com/delaneyj/datastar"
	"github.com/gorilla/sessions"
	"github.com/ngsalvo/roadmapsh-personal-blog/internal"
)

type PostLogout struct {
	session sessions.Store
}

func NewPostLogout(session sessions.Store) PostLogout {
	return PostLogout{
		session: session,
	}
}

func (h PostLogout) Handle(w http.ResponseWriter, r *http.Request) {
	internal.ClearSession(h.session, w, r)

	sse := datastar.NewSSE(w, r)

	datastar.Redirect(sse, "/")

}
