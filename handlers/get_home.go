package handlers

import (
	"log"
	"net/http"
	"strings"

	"github.com/adrg/frontmatter"
	"github.com/ngsalvo/roadmapsh-personal-blog/components"
	"github.com/ngsalvo/roadmapsh-personal-blog/dtos"
	"github.com/ngsalvo/roadmapsh-personal-blog/services"
)

type GetHome interface {
	Handle(w http.ResponseWriter, r *http.Request)
}

type getHome struct {
	fileReader services.FileReader
}

func NewGetHome(fileReader services.FileReader) GetHome {
	return &getHome{
		fileReader: fileReader,
	}
}

func (h *getHome) Handle(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html")

	slugs, err := h.fileReader.GetFileNames("static/blog")
	if err != nil {
		http.Error(w, "Error reading blog directory", http.StatusInternalServerError)
		return
	}

	articles := make([]dtos.Article, len(slugs))

	for i, fileName := range slugs {
		article, err := h.fileReader.Read("static/blog/" + fileName)
		if err != nil {
			http.Error(w, "Article not found", http.StatusNotFound)
			return
		}

		var articleData dtos.Article

		_, err = frontmatter.Parse(strings.NewReader(article), &articleData)

		if err != nil {
			http.Error(w, "Error parsing frontmatter", http.StatusInternalServerError)
			log.Fatal(err)
		}

		articleData.Slug = slugs[i]
		articles[i] = articleData
	}

	components.Home(articles).Render(r.Context(), w)
}
