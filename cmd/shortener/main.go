package main

import (
	"net/http"
	"sprint/cmd/shortener/config"
	"sprint/cmd/shortener/handlers"
	"strconv"

	"github.com/go-chi/chi/v5"
)

func main() {
	if err := run(); err != nil {
		panic(err)
	}
}

func run() error {
	flags := config.ParseFlags()
	r := chi.NewRouter()
	r.Route("/", func(r chi.Router) {
		r.Post("/", func(w http.ResponseWriter, r *http.Request) {
			handlers.HandlerPost(w, r, flags)
		})
	})
	r.Route("/{id:[a-zA-Z0-9]+}", func(r chi.Router) {
		r.Get("/", handlers.HandlerGet)
	})
	// fmt.Println(flags.Host+":"+strconv.Itoa(flags.Port), flags.Port, flags.BaseURL)
	adress := flags.Host + ":" + strconv.Itoa(flags.Port)
	return http.ListenAndServe(adress, r)
}
