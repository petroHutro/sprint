package router

import (
	"net/http"
	"sprint/internal/config"
	"sprint/internal/handlers"
	"sprint/internal/logger"

	"github.com/go-chi/chi/v5"
	"go.uber.org/zap"
)

func Router(flags *config.Flags, log *zap.Logger) *chi.Mux {
	r := chi.NewRouter()
	r.Use(logger.LoggingMiddleware(log))
	r.Route("/", func(r chi.Router) {
		r.Post("/", func(w http.ResponseWriter, r *http.Request) {
			handlers.HandlerPost(w, r, string(flags.BaseURL))
		})
	})
	r.Route("/{id:[a-zA-Z0-9]+}", func(r chi.Router) {
		r.Get("/", handlers.HandlerGet)
	})
	return r
}
