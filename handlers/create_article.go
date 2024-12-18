package handlers

import (
	"log"
	"net/http"

	"github.com/ngsalvo/roadmapsh-personal-blog/dtos"
	datastar "github.com/starfederation/datastar/sdk/go"

	"github.com/ngsalvo/roadmapsh-personal-blog/repositories"
)

type CreateArticle interface {
	Handle(w http.ResponseWriter, r *http.Request)
}

type createArticle struct {
	fileReader repositories.FileReader
}

func NewCreateArticle(fileReader repositories.FileReader) CreateArticle {
	return &createArticle{
		fileReader: fileReader,
	}
}

func (h *createArticle) Handle(w http.ResponseWriter, r *http.Request) {

	var article dtos.NewArticle
	err := datastar.ReadSignals(r, &article)

	if err != nil {
		log.Println(err)
		return
	}

	err = h.fileReader.Create(&article)

	if err != nil {
		log.Println(err)
		return
	}

	sse := datastar.NewSSE(w, r)
	sse.Redirect("/admin")

	return
}
