package main

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/gorilla/sessions"
	"github.com/ngsalvo/roadmapsh-personal-blog/components"
	"github.com/ngsalvo/roadmapsh-personal-blog/datasources"
	"github.com/ngsalvo/roadmapsh-personal-blog/handlers"
	"github.com/ngsalvo/roadmapsh-personal-blog/middlewares"
	"github.com/ngsalvo/roadmapsh-personal-blog/repositories"
)

func ConfigureRoutes(r chi.Router) {
	fileReader := repositories.NewFileReader()
	articleDatasource := datasources.NewArticleDatasource(fileReader)
	fileServer := http.FileServer(http.Dir("./static"))
	session := sessions.NewCookieStore([]byte("2La74FjQAk9N7zrBwO9hmL7qgzYQFeVx"))
	session.Options = &sessions.Options{
		HttpOnly: true,
		SameSite: http.SameSiteLaxMode,
	}
	authenthication := middlewares.NewAuthMiddleware(session)

	r.Handle("/static/*", http.StripPrefix("/static/", fileServer))
	r.Handle("/", http.RedirectHandler("/home", http.StatusPermanentRedirect))

	r.Get("/home", handlers.NewGetHome(articleDatasource, session).Handle)

	r.Get("/register", func(w http.ResponseWriter, r *http.Request) {
		components.SignIn().Render(r.Context(), w)
	})
	r.Post("/register", handlers.NewPostRegister(session).Handle)
	r.Get("/login", handlers.NewGetLogin().Handle)
	r.Post("/login", handlers.NewPostLogin(session).Handle)
	r.Post("/logout", handlers.NewPostLogout(session).Handle)

	r.Get("/article/{slug}", handlers.NewGetArticle(articleDatasource).Handle)

	r.Group(func(r chi.Router) {
		r.Use(authenthication.Authentication)
		r.Get("/admin", handlers.NewGetAdmin(articleDatasource).Handle)
		r.Get("/article/new", handlers.NewGetCreateArticle().Handle)
		r.Post("/article/new", handlers.NewCreateArticle(fileReader).Handle)
		r.Get("/article/{slug}/edit", handlers.NewGetArticleEdit(articleDatasource).Handle)
		r.Put("/article/{slug}/edit", handlers.NewUpdateArticle(fileReader).Handle)
		r.Delete("/article/{slug}/delete", handlers.NewDeleteArticle(fileReader).Handle)
	})
}
