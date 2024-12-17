package handlers

import (
	"net/http"

	"github.com/ngsalvo/roadmapsh-personal-blog/components"
)

type GetRegister interface {
	Handle(w http.ResponseWriter, r *http.Request)
}

type getRegister struct {
}

func NewGetRegister() GetLogin {
	return &getRegister{}
}

func (h *getRegister) Handle(w http.ResponseWriter, r *http.Request) {
	components.SignIn().Render(r.Context(), w)
}
