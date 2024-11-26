package handlers

import (
	"net/http"

	"github.com/ngsalvo/roadmapsh-personal-blog/components"
)

type GetArticleNew interface {
	Handle(w http.ResponseWriter, r *http.Request)
}

type getArticleNew struct {
}

func NewGetCreateArticle() GetArticleNew {
	return &getArticleNew{}
}

func (h *getArticleNew) Handle(w http.ResponseWriter, r *http.Request) {
	components.NewArticle().Render(r.Context(), w)
}
