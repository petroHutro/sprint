package router

import (
	"net/http"
	"sprint/internal/authorization"
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
	r.Use(authorization.AuthorizationMiddleware)
	r.Use(compression.GzipMiddleware)

	r.Route("/", func(r chi.Router) {
		r.Post("/", func(w http.ResponseWriter, r *http.Request) {
			handlers.HandlerPost(w, r, flags.BaseURL, flags.FileStoragePath, db)
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

	r.Route("/api", func(r chi.Router) {
		r.Post("/shorten", func(w http.ResponseWriter, r *http.Request) {
			handlers.HandlerPostAPI(w, r, flags.BaseURL, flags.FileStoragePath, db)
		})
		r.Post("/shorten/batch", func(w http.ResponseWriter, r *http.Request) {
			handlers.HandlerPostBatch(w, r, flags.BaseURL, flags.FileStoragePath, db)
		})
		r.Get("/user/urls", func(w http.ResponseWriter, r *http.Request) {
			handlers.HandlerGetUrls(w, r, flags.BaseURL, db)
		})
	})

	return r
}
