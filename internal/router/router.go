package router

import (
	"net/http"
	"sprint/internal/config"
	"sprint/internal/handlers"

	"github.com/go-chi/chi/v5"
)

func Router(flags *config.Flags) *chi.Mux {
	r := chi.NewRouter()
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
