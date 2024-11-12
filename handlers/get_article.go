package handlers

import (
	"bytes"
	"context"
	"errors"
	"io"
	"log"
	"net/http"
	"strings"

	"github.com/a-h/templ"
	"github.com/adrg/frontmatter"
	"github.com/go-chi/chi/v5"
	"github.com/ngsalvo/roadmapsh-personal-blog/components"
	"github.com/ngsalvo/roadmapsh-personal-blog/dtos"
	customErrors "github.com/ngsalvo/roadmapsh-personal-blog/errors"
	"github.com/ngsalvo/roadmapsh-personal-blog/repositories"
	"github.com/yuin/goldmark"
	highlighting "github.com/yuin/goldmark-highlighting/v2"
)

type GetArticle interface {
	Handle(w http.ResponseWriter, r *http.Request)
}

type getArticle struct {
	fileReader repositories.FileReader
}

func NewGetArticle(fileReader repositories.FileReader) GetArticle {
	return &getArticle{
		fileReader: fileReader,
	}
}

func (h *getArticle) Handle(w http.ResponseWriter, r *http.Request) {
	slug := chi.URLParam(r, "slug")

	article, err := h.fileReader.Read("static/blog/" + slug)
	if err != nil {
		var applicationError customErrors.ApplicationError
		if errors.As(err, &applicationError) {
			if strings.Contains(applicationError.Message, "article not found") {
				http.Error(w, "Article not found", http.StatusNotFound)
				return
			}
		}
		http.Error(w, "Error reading file", http.StatusInternalServerError)
		return
	}

	var articleData dtos.Article

	remainingMD, err := frontmatter.Parse(strings.NewReader(article), &articleData)

	if err != nil {
		http.Error(w, "Error parsing frontmatter", http.StatusInternalServerError)
		log.Fatal(err)
	}

	mdRenderer := goldmark.New(
		goldmark.WithExtensions(
			highlighting.NewHighlighting(
				highlighting.WithStyle("tokyonight-storm")),
		),
	)
	var buffer bytes.Buffer
	err = mdRenderer.Convert([]byte(remainingMD), &buffer)
	if err != nil {
		http.Error(w, "Error reading file", http.StatusInternalServerError)
		log.Fatal(err)
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

	content := unsafe(buffer.String())

	articleData.Content = content

	components.Article(articleData).Render(r.Context(), w)
}
