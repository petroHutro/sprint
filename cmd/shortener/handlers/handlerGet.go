package handlers

import (
	"net/http"
	"sprint/cmd/shortener/db"
)

func HandlerGet(w http.ResponseWriter, r *http.Request) {
	shortLink := r.URL.Path[1:]
	if shortLink == "" {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	var link string
	for key, value := range db.Db {
		if value == shortLink {
			link = key
		}
	}
	if link == "" {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	w.Header().Set("Location", link)
	w.WriteHeader(http.StatusTemporaryRedirect)
	w.Write([]byte{})
}
