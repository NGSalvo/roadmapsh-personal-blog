package handlers

import (
	"log"
	"net/http"

	"github.com/delaneyj/datastar"
	"github.com/gorilla/sessions"
	"github.com/ngsalvo/roadmapsh-personal-blog/dtos"
	"github.com/ngsalvo/roadmapsh-personal-blog/internal"
)

type PostRegister interface {
	Handle(w http.ResponseWriter, r *http.Request)
}

type postRegister struct {
	session sessions.Store
}

func NewPostRegister(session sessions.Store) PostRegister {
	return postRegister{
		session: session,
	}
}

func (pr postRegister) Handle(w http.ResponseWriter, r *http.Request) {
	var store dtos.UserLogin
	err := datastar.BodyUnmarshal(r, &store)

	if err != nil {
		log.Println(err)
		return
	}

	hashPassword, err := internal.HashPassword(store.Password)

	if err != nil {
		log.Println(err)
		return
	}

	internal.AddUser(store.Username, hashPassword)
	internal.CreateSession(pr.session, w, r, store.Username)

	// err = repositories.CreateUser(&store)

	sse := datastar.NewSSE(w, r)

	datastar.Redirect(sse, "/")
}
