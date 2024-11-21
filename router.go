package main

import (
	"log"
	"net/http"

	"github.com/delaneyj/datastar"
	"github.com/go-chi/chi/v5"
	"github.com/ngsalvo/roadmapsh-personal-blog/datasources"
	"github.com/ngsalvo/roadmapsh-personal-blog/dtos"
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

	r.Get("/article/{slug}", handlers.NewGetArticle(articleDatasource).Handle)
	r.Get("/article/{slug}/edit", handlers.NewGetArticleEdit(articleDatasource).Handle)
	r.Put("/article/{slug}/edit", func(w http.ResponseWriter, r *http.Request) {
		log.Println("PUT /article/{slug}/edit")

		slug := chi.URLParam(r, "slug")

		log.Printf("slug: %s", slug)

		var store dtos.ArticleStore
		err := datastar.BodyUnmarshal(r, &store)

		if err != nil {
			log.Println(err)
			return
		}

		fr := repositories.NewFileReader()

		err = fr.Update(slug, &store)

		if err != nil {
			log.Println(err)
			return
		}

		sse := datastar.NewSSE(w, r)
		datastar.Redirect(sse, "/admin")

		return
	})

}
