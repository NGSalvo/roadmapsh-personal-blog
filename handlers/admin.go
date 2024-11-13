package handlers

import (
	"errors"
	"net/http"
	"strings"

	"github.com/ngsalvo/roadmapsh-personal-blog/components"
	"github.com/ngsalvo/roadmapsh-personal-blog/datasources"
	customErrors "github.com/ngsalvo/roadmapsh-personal-blog/errors"
)

type (
	GetAdmin interface {
		Handle(w http.ResponseWriter, r *http.Request)
	}

	getAdmin struct {
		articleDatasource datasources.ArticlesDatasource
	}
)

func NewGetAdmin(articleDatasource datasources.ArticlesDatasource) GetAdmin {
	return &getAdmin{
		articleDatasource: articleDatasource,
	}
}

func (h *getAdmin) Handle(w http.ResponseWriter, r *http.Request) {
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

	components.Dashboard(articles).Render(r.Context(), w)
}
