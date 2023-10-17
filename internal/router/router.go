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

func Create(flags *config.Flags, db *storage.StorageBase) *chi.Mux {
	r := chi.NewRouter()
	r.Use(logger.LoggingMiddleware)
	r.Use(compression.GzipMiddleware)

	r.Route("/", func(r chi.Router) {
		r.Post("/", func(w http.ResponseWriter, r *http.Request) {
			handlers.HandlerPost(w, r, string(flags.BaseURL), flags.FileStoragePath, db)
		})
	})

	r.Route("/ping", func(r chi.Router) {
		r.Get("/", func(w http.ResponseWriter, r *http.Request) {
			handlers.HandlerPing(w, r, db)
		})
	})

	r.Route("/{id:[a-zA-Z0-9]+}", func(r chi.Router) {
		r.Get("/", func(w http.ResponseWriter, r *http.Request) {
			handlers.HandlerGet(w, r, db)
		})
	})

	r.Route("/api/shorten", func(r chi.Router) {
		r.Post("/", func(w http.ResponseWriter, r *http.Request) {
			handlers.HandlerPostAPI(w, r, string(flags.BaseURL), flags.FileStoragePath, db)
		})
		r.Post("/batch", func(w http.ResponseWriter, r *http.Request) {
			handlers.HandlerPostBatch(w, r, string(flags.BaseURL), flags.FileStoragePath, db)
		})
	})

	return r
}
