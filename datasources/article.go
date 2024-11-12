package datasources

import (
	"errors"
	"fmt"

	"github.com/ngsalvo/roadmapsh-personal-blog/dtos"
	customErrors "github.com/ngsalvo/roadmapsh-personal-blog/errors"
	"github.com/ngsalvo/roadmapsh-personal-blog/repositories"
	"github.com/ngsalvo/roadmapsh-personal-blog/services"
)

type (
	ArticlesDatasource interface {
		GetArticles() ([]dtos.Article, error)
	}

	articlesDatasource struct {
		fileReader repositories.FileReader
	}
)

const (
	articleNotFound          = "article not found"
	articleDirectoryNotFound = "article directory not found"
)

func NewArticleDatasource(fileReader repositories.FileReader) ArticlesDatasource {
	return &articlesDatasource{
		fileReader: fileReader,
	}
}

func (s *articlesDatasource) GetArticles() ([]dtos.Article, error) {
	slugs, err := s.fileReader.GetFileNames("static/blog")

	if err != nil {
		if errors.Is(err, customErrors.ErrorArticleNotFound) {
			return nil, &customErrors.ApplicationError{
				Message: fmt.Sprintf("%s: %s", articleNotFound, err.Error()),
			}
		}
		if errors.Is(err, customErrors.ErrorReadingDirectory) {
			return nil, &customErrors.ApplicationError{
				Message: fmt.Sprintf("%s: %s", articleDirectoryNotFound, err.Error()),
			}
		}
		return nil, fmt.Errorf("error getting all titles: %w", err)
	}

	articles := make([]dtos.Article, len(slugs))
	var articleData *dtos.Article

	for i, fileName := range slugs {
		article, err := s.fileReader.Read("static/blog/" + fileName)
		if err != nil {
			if errors.Is(err, customErrors.ErrorArticleNotFound) {
				return nil, &customErrors.ApplicationError{
					Message: fmt.Sprintf("%s: %s", articleNotFound, fileName),
				}
			}
		}

		articleData, err = services.Parse[dtos.Article](article)

		if err != nil {
			return nil, err
		}

		articleData.Slug = slugs[i]
		articles[i] = *articleData
	}

	return articles, nil
}
