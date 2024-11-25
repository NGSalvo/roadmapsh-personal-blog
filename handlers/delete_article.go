package handlers

import (
	"log"
	"net/http"

	"github.com/delaneyj/datastar"
	"github.com/go-chi/chi/v5"
	"github.com/ngsalvo/roadmapsh-personal-blog/dtos"

	"github.com/ngsalvo/roadmapsh-personal-blog/repositories"
)

type DeleteArticle interface {
	Handle(w http.ResponseWriter, r *http.Request)
}

type deleteArticle struct {
	fileReader repositories.FileReader
}

func NewDeleteArticle(fileReader repositories.FileReader) DeleteArticle {
	return &deleteArticle{
		fileReader: fileReader,
	}
}

func (h *deleteArticle) Handle(w http.ResponseWriter, r *http.Request) {
	slug := chi.URLParam(r, "slug")

	var store dtos.ArticleStore
	err := datastar.BodyUnmarshal(r, &store)

	if err != nil {
		log.Println(err)
		return
	}

	err = h.fileReader.Delete(slug)

	if err != nil {
		log.Println(err)
		return
	}

	sse := datastar.NewSSE(w, r)
	datastar.Redirect(sse, "/admin")

	return
}
