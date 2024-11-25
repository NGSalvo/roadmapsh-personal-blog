package handlers

import (
	"log"
	"net/http"

	"github.com/delaneyj/datastar"
	"github.com/ngsalvo/roadmapsh-personal-blog/dtos"

	"github.com/go-chi/chi/v5"
	"github.com/ngsalvo/roadmapsh-personal-blog/repositories"
)

type UpdateArticle interface {
	Handle(w http.ResponseWriter, r *http.Request)
}

type updateArticle struct {
	fileReader repositories.FileReader
}

func NewUpdateArticle(fileReader repositories.FileReader) UpdateArticle {
	return &updateArticle{
		fileReader: fileReader,
	}
}

func (h *updateArticle) Handle(w http.ResponseWriter, r *http.Request) {
	slug := chi.URLParam(r, "slug")

	var store dtos.ArticleStore
	err := datastar.BodyUnmarshal(r, &store)

	if err != nil {
		log.Println(err)
		return
	}

	fr := repositories.NewFileReader()

	err = fr.Update(slug, &store)

	if err != nil {
		log.Println(err)
		return
	}

	sse := datastar.NewSSE(w, r)
	datastar.Redirect(sse, "/admin")

	return
}
