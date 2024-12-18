package handlers

import (
	"log"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/ngsalvo/roadmapsh-personal-blog/dtos"
	datastar "github.com/starfederation/datastar/sdk/go"

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
	err := datastar.ReadSignals(r, &store)

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
	sse.Redirect("/admin")

	return
}
