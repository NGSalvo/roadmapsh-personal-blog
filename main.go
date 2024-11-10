package main

import (
	"log"
	"log/slog"
	"net/http"
	"os"

	"github.com/go-chi/chi/v5"
)

func main() {
	logger := slog.New(slog.NewJSONHandler(os.Stdout, nil))

	router := chi.NewRouter()

	ConfigureRoutes(router)

	logger.Info("Starting server on port 3000")
	log.Fatal(http.ListenAndServe(":3000", router))
}
