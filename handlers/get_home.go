package handlers

import (
	"errors"
	"net/http"
	"strings"

	"github.com/ngsalvo/roadmapsh-personal-blog/components"
	"github.com/ngsalvo/roadmapsh-personal-blog/datasources"
	customErrors "github.com/ngsalvo/roadmapsh-personal-blog/errors"
)

type GetHome interface {
	Handle(w http.ResponseWriter, r *http.Request)
}

type getHome struct {
	articleDatasource datasources.ArticlesDatasource
}

func NewGetHome(articleDatasource datasources.ArticlesDatasource) GetHome {
	return &getHome{
		articleDatasource: articleDatasource,
	}
}

func (h *getHome) Handle(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html")

	articles, err := h.articleDatasource.GetArticles()

	if err != nil {
		var applicationError customErrors.ApplicationError
		if errors.As(err, &applicationError) {
			if strings.Contains(applicationError.Message, "article not found") {
				http.Error(w, "Article not found", http.StatusNotFound)
				return
			}

			if strings.Contains(applicationError.Message, "article directory not found") {
				http.Error(w, "Article directory not found", http.StatusNotFound)
				return
			}
		}

		http.Error(w, "Error getting articles", http.StatusInternalServerError)
		return
	}

	components.Home(articles).Render(r.Context(), w)
}
