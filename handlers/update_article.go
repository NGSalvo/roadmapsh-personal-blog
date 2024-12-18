package handlers

import (
	"log"
	"net/http"

	"github.com/ngsalvo/roadmapsh-personal-blog/dtos"
	datastar "github.com/starfederation/datastar/sdk/go"

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
	err := datastar.ReadSignals(r, &store)

	if err != nil {
		log.Println(err)
		return
	}

	err = h.fileReader.Update(slug, &store)

	if err != nil {
		log.Println(err)
		return
	}

	sse := datastar.NewSSE(w, r)
	sse.Redirect("/admin")

	return
}
