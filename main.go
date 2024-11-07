package main

import (
	"bytes"
	"context"
	"io"
	"log"
	"log/slog"
	"net/http"
	"os"
	"strings"

	"github.com/a-h/templ"
	"github.com/adrg/frontmatter"
	"github.com/go-chi/chi/v5"
	"github.com/yuin/goldmark"
	highlighting "github.com/yuin/goldmark-highlighting/v2"

	"github.com/ngsalvo/roadmapsh-personal-blog/components"
)

func main() {
	fileServer := http.FileServer(http.Dir("./static"))
	logger := slog.New(slog.NewJSONHandler(os.Stdout, nil))
	fileReader := NewFileReader()
	getArticleHandler := NewGetArticle(fileReader).Handle

	router := chi.NewRouter()

	router.Handle("/static/*", http.StripPrefix("/static/", fileServer))
	router.Handle("/", http.RedirectHandler("/home", http.StatusPermanentRedirect))
	router.Get("/home", homeHandler)
	router.Get("/article/{slug}", getArticleHandler)

	logger.Info("Starting server on port 3000")
	log.Fatal(http.ListenAndServe(":3000", router))
}

func homeHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html")

	readFromAssets := func(path string) ([]string, error) {
		dir, err := os.ReadDir(path)
		if err != nil {
			return nil, err
		}

		var files []string
		for _, file := range dir {
			fileName := strings.Split(file.Name(), ".")[0]
			files = append(files, fileName)
		}

		return files, nil
	}

	fileReader := NewFileReader()

	slugs, err := readFromAssets("static/blog")
	if err != nil {
		http.Error(w, "Error reading blog directory", http.StatusInternalServerError)
		return
	}

	articles := make([]components.ArticleStore, len(slugs))

	for i, fileName := range slugs {
		article, err := fileReader.Read("static/blog/" + fileName)
		if err != nil {
			http.Error(w, "Article not found", http.StatusNotFound)
			return
		}

		var articleData components.ArticleStore

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

func getArticleHandler(w http.ResponseWriter, r *http.Request) {
	slug := chi.URLParam(r, "slug")
	fileReader := NewFileReader()

	logger := slog.New(slog.NewJSONHandler(os.Stdout, nil))
	article, err := fileReader.Read("static/blog/" + slug)
	if err != nil {
		http.Error(w, "Article not found", http.StatusNotFound)
		return
	}

	var articleData components.ArticleStore

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
	logger.Info(buffer.String())

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

type FileReader struct{}

type GetArticle interface {
	Handle(w http.ResponseWriter, r *http.Request)
}

type getArticle struct {
	fileReader FileReader
}

func NewGetArticle(fileReader FileReader) GetArticle {
	return &getArticle{
		fileReader: fileReader,
	}
}

func (g *getArticle) Handle(w http.ResponseWriter, r *http.Request) {
	getArticleHandler(w, r)
}

func NewFileReader() FileReader {
	return FileReader{}
}

func (fr FileReader) Read(slug string) (string, error) {
	f, err := os.Open(slug + ".md")
	if err != nil {
		return "", nil
	}

	defer f.Close()

	bytes, err := io.ReadAll(f)
	if err != nil {
		return "", nil
	}

	return string(bytes), nil
}
