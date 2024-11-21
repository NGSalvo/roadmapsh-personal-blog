package datasources

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"io"

	"github.com/a-h/templ"
	"github.com/ngsalvo/roadmapsh-personal-blog/dtos"
	customErrors "github.com/ngsalvo/roadmapsh-personal-blog/errors"
	"github.com/ngsalvo/roadmapsh-personal-blog/repositories"
	"github.com/ngsalvo/roadmapsh-personal-blog/services"
	"github.com/yuin/goldmark"
	highlighting "github.com/yuin/goldmark-highlighting/v2"
)

type (
	ArticlesDatasource interface {
		GetArticles() ([]dtos.Article, error)
		GetArticle(slug string) (*dtos.Article, error)
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

func (d *articlesDatasource) GetArticles() ([]dtos.Article, error) {
	slugs, err := d.fileReader.GetFileNames("static/blog")

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
	var articleData *services.Frontmatter[dtos.Article]

	for i, fileName := range slugs {
		article, err := d.fileReader.Read("static/blog/" + fileName)
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

		articleData.Frontmatter.Slug = slugs[i]
		articles[i] = *&articleData.Frontmatter
	}

	return articles, nil
}

func (d *articlesDatasource) GetArticle(slug string) (*dtos.Article, error) {
	article, err := d.fileReader.Read("static/blog/" + slug)
	if err != nil {
		if errors.Is(err, customErrors.ErrorArticleNotFound) {
			return nil, &customErrors.ApplicationError{
				Message: fmt.Sprintf("%s: %s", articleNotFound, err.Error()),
			}
		}
		return nil, fmt.Errorf("error getting article: %w", err)
	}

	var articleData *services.Frontmatter[dtos.Article]

	articleData, err = services.Parse[dtos.Article](article)

	if err != nil {
		return nil, fmt.Errorf("error parsing frontmatter: %w", err)
	}

	mdRenderer := goldmark.New(
		goldmark.WithExtensions(
			highlighting.NewHighlighting(
				highlighting.WithStyle("tokyonight-storm")),
		),
	)
	var buffer bytes.Buffer
	err = mdRenderer.Convert([]byte(articleData.RemaingData), &buffer)
	if err != nil {
		return nil, fmt.Errorf("error rendering markdown: %w", err)
	}

	unsafe := func(html string) templ.Component {
		return templ.ComponentFunc(func(ctx context.Context, w io.Writer) error {
			_, err := io.WriteString(w, html)
			if err != nil {
				panic(err)
			}
			return nil
		})
	}

	contentString := string(articleData.RemaingData)
	articleData.Frontmatter.ContentString = contentString

	content := unsafe(buffer.String())

	articleData.Frontmatter.Content = content
	articleData.Frontmatter.Slug = slug

	return &articleData.Frontmatter, nil
}
