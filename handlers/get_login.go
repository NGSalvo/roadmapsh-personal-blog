package handlers

import (
	"net/http"

	"github.com/ngsalvo/roadmapsh-personal-blog/components"
	"github.com/ngsalvo/roadmapsh-personal-blog/dtos"
)

type GetLogin interface {
	Handle(w http.ResponseWriter, r *http.Request)
}

type getLogin struct {
}

func NewGetLogin() GetLogin {
	return &getLogin{}
}

func (h *getLogin) Handle(w http.ResponseWriter, r *http.Request) {
	components.Login(dtos.UserLogin{}).Render(r.Context(), w)
}
