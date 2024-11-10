package main

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/ngsalvo/roadmapsh-personal-blog/handlers"
	"github.com/ngsalvo/roadmapsh-personal-blog/services"
)

func ConfigureRoutes(r chi.Router) {
	fileReader := services.NewFileReader()
	fileServer := http.FileServer(http.Dir("./static"))

	r.Handle("/static/*", http.StripPrefix("/static/", fileServer))
	r.Handle("/", http.RedirectHandler("/home", http.StatusPermanentRedirect))

	r.Get("/home", handlers.NewGetHome(fileReader).Handle)
	r.Get("/article/{slug}", handlers.NewGetArticle(fileReader).Handle)
}
