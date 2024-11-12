package handlers

import (
	"log"
	"net/http"

	"github.com/ngsalvo/roadmapsh-personal-blog/components"
	"github.com/ngsalvo/roadmapsh-personal-blog/dtos"
	"github.com/ngsalvo/roadmapsh-personal-blog/repositories"
	"github.com/ngsalvo/roadmapsh-personal-blog/services"
)

type (
	GetAdmin interface {
		Handle(w http.ResponseWriter, r *http.Request)
	}

	getAdmin struct {
		fileReader repositories.FileReader
	}
)

func NewGetAdmin(fileReader repositories.FileReader) GetAdmin {
	return &getAdmin{
		fileReader: fileReader,
	}
}

func (h *getAdmin) Handle(w http.ResponseWriter, r *http.Request) {
	slugs, err := h.fileReader.GetFileNames("static/blog")

	if err != nil {
		http.Error(w, "Error reading blog directory", http.StatusInternalServerError)
		return
	}

	articles := make([]dtos.Article, len(slugs))
	var articleData *dtos.Article

	for i, fileName := range slugs {
		article, err := h.fileReader.Read("static/blog/" + fileName)
		if err != nil {
			http.Error(w, "Article not found", http.StatusNotFound)
			return
		}

		articleData, err = services.Parse[dtos.Article](article)

		if err != nil {
			http.Error(w, "Error parsing frontmatter", http.StatusInternalServerError)
			log.Fatal(err)
		}

		articleData.Slug = slugs[i]
		articles[i] = *articleData
	}

	components.Dashboard(articles).Render(r.Context(), w)
}
