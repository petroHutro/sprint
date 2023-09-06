package main

import (
	"net/http"
	"sprint/cmd/shortener/handlers"

	"github.com/gorilla/mux"
)

func main() {
	if err := run(); err != nil {
		panic(err)
	}
}

func run() error {
	r := mux.NewRouter()
	r.HandleFunc("/", handlers.HandlerPost).Methods("POST")
	r.HandleFunc("/{id:[a-zA-Z0-9]+}", handlers.HandlerGet).Methods("GET")
	return http.ListenAndServe(`:8080`, r)
}
