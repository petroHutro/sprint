package main

import (
	"net/http"
	"sprint/cmd/shortener/handlers"

	// "github.com/gorilla/mux"
	"github.com/go-chi/chi/v5"
)

func main() {
	if err := run(); err != nil {
		panic(err)
	}
}

func run() error {
	r := chi.NewRouter()

	r.Route("/", func(r chi.Router) {
		r.Post("/", handlers.HandlerPost)
	})
	r.Route("/{id:[a-zA-Z0-9]+}", func(r chi.Router) {
		r.Get("/", handlers.HandlerGet)
	})

	// r := mux.NewRouter()
	// r.HandleFunc("/", handlers.HandlerPost).Methods("POST")
	// r.HandleFunc("/{id:[a-zA-Z0-9]+}", handlers.HandlerGet).Methods("GET")
	return http.ListenAndServe(`:8080`, r)
}
