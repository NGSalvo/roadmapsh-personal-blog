package handlers

import (
	"log"
	"net/http"

	"github.com/delaneyj/datastar"
	"github.com/gorilla/sessions"

	"github.com/ngsalvo/roadmapsh-personal-blog/dtos"
	"github.com/ngsalvo/roadmapsh-personal-blog/internal"
	auth "github.com/ngsalvo/roadmapsh-personal-blog/internal"
)

type PostLogin interface {
	Handle(w http.ResponseWriter, r *http.Request)
}

type postLogin struct {
	session sessions.Store
}

func NewPostLogin(session sessions.Store) PostLogin {
	return &postLogin{
		session: session,
	}
}

type contextKey string

const UserKey contextKey = "username"

func (h *postLogin) Handle(w http.ResponseWriter, r *http.Request) {
	log.Println("POST /login")

	var store dtos.UserLogin
	err := datastar.BodyUnmarshal(r, &store)

	if err != nil {
		log.Println("Unmarshal error: ", err)
		return
	}

	user := auth.GetUser(store.Username)

	if user == nil || !auth.CheckPasswordHash(store.Password, user.HashedPassword) {
		sse := datastar.NewSSE(w, r)

		datastar.PatchStore(sse, dtos.UserLogin{
			Username: store.Username,
			Password: "",
		})

		datastar.RenderFragmentString(sse, "<div id=\"login-error\" class=\"alert alert-danger\">Invalid username or password</div>", datastar.WithQuerySelectorID("login-error"))

		return
	}

	internal.CreateSession(h.session, w, r, store.Username)

	if err != nil {
		log.Println(err)
		return
	}

	sse := datastar.NewSSE(w, r)

	datastar.Redirect(sse, "/")
}
