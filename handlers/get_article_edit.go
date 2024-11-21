package handlers

import (
	"errors"
	"net/http"
	"strings"

	"github.com/go-chi/chi/v5"
	"github.com/ngsalvo/roadmapsh-personal-blog/components"
	"github.com/ngsalvo/roadmapsh-personal-blog/datasources"
	customErrors "github.com/ngsalvo/roadmapsh-personal-blog/errors"
)

type GetArticleEdit interface {
	Handle(w http.ResponseWriter, r *http.Request)
}

type getArticleEdit struct {
	articleDatasource datasources.ArticlesDatasource
}

func NewGetArticleEdit(articleDatasource datasources.ArticlesDatasource) GetArticleEdit {
	return &getArticleEdit{
		articleDatasource: articleDatasource,
	}
}

func (h *getArticleEdit) Handle(w http.ResponseWriter, r *http.Request) {
	slug := chi.URLParam(r, "slug")

	article, err := h.articleDatasource.GetArticle(slug)

	if err != nil {
		var applicationError customErrors.ApplicationError
		if errors.As(err, &applicationError) {
			if strings.Contains(applicationError.Message, "article not found") {
				http.Error(w, "Article not found", http.StatusNotFound)
				return
			}
		}

		http.Error(w, "Error getting article", http.StatusInternalServerError)
		return
	}

	components.ArticleEdit(article).Render(r.Context(), w)
}
