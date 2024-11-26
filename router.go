package main

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/ngsalvo/roadmapsh-personal-blog/datasources"
	"github.com/ngsalvo/roadmapsh-personal-blog/handlers"
	"github.com/ngsalvo/roadmapsh-personal-blog/repositories"
)

func ConfigureRoutes(r chi.Router) {
	fileReader := repositories.NewFileReader()
	articleDatasource := datasources.NewArticleDatasource(fileReader)
	fileServer := http.FileServer(http.Dir("./static"))

	r.Handle("/static/*", http.StripPrefix("/static/", fileServer))
	r.Handle("/", http.RedirectHandler("/home", http.StatusPermanentRedirect))

	r.Get("/home", handlers.NewGetHome(articleDatasource).Handle)
	r.Get("/admin", handlers.NewGetAdmin(articleDatasource).Handle)

	r.Get("/article/new", handlers.NewGetCreateArticle().Handle)
	r.Post("/article/new", handlers.NewCreateArticle(fileReader).Handle)
	r.Get("/article/{slug}", handlers.NewGetArticle(articleDatasource).Handle)
	r.Get("/article/{slug}/edit", handlers.NewGetArticleEdit(articleDatasource).Handle)
	r.Put("/article/{slug}/edit", handlers.NewUpdateArticle(fileReader).Handle)
	r.Delete("/article/{slug}/delete", handlers.NewDeleteArticle(fileReader).Handle)
}
