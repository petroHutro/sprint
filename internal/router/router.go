package router

import (
	"net/http"
	"sprint/internal/compression"
	"sprint/internal/config"
	"sprint/internal/handlers"
	"sprint/internal/logger"
	"sprint/internal/storage"

	"github.com/go-chi/chi/v5"
)

func Create(flags *config.Flags) *chi.Mux {
	db := storage.Connection(string(flags.DatabaseDSN))
	r := chi.NewRouter()
	r.Use(logger.LoggingMiddleware)
	r.Use(compression.GzipMiddleware)
	r.Route("/", func(r chi.Router) {
		r.Post("/", func(w http.ResponseWriter, r *http.Request) {
			handlers.HandlerPost(w, r, string(flags.BaseURL), string(flags.FileStoragePath))
		})
	})
	r.Route("/ping", func(r chi.Router) {
		r.Get("/", func(w http.ResponseWriter, r *http.Request) {
			handlers.HandlerPing(w, r, db)
		})
	})
	r.Route("/{id:[a-zA-Z0-9]+}", func(r chi.Router) {
		r.Get("/", handlers.HandlerGet)
	})
	r.Route("/api", func(r chi.Router) {
		r.Post("/shorten", func(w http.ResponseWriter, r *http.Request) {
			handlers.HandlerPostAPI(w, r, string(flags.BaseURL), string(flags.FileStoragePath))
		})
	})
	return r
}
