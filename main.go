package main

import (
	"log"
	"log/slog"
	"net/http"
	"os"

	"github.com/go-chi/chi/v5"

	"github.com/ngsalvo/roadmapsh-personal-blog/components"
)

func main() {
	fileServer := http.FileServer(http.Dir("./static"))
	logger := slog.New(slog.NewJSONHandler(os.Stdout, nil))

	router := chi.NewRouter()

	router.Get("/", homeHandler)
	router.Handle("/static/*", http.StripPrefix("/static/", fileServer))

	logger.Info("Starting server on port 3000")
	log.Fatal(http.ListenAndServe(":3000", router))
}

func homeHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html")

	components.Home().Render(r.Context(), w)
}
