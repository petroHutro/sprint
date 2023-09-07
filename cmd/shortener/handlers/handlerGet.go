package handlers

import (
	"net/http"
	"sprint/cmd/shortener/db"
)

func HandlerGet(w http.ResponseWriter, r *http.Request) {
	shortLink := r.URL.String()[1:]
	// shortLink := r.URL.Query().Get("id")
	// shortLink := chi.URLParam(r, "id")
	if shortLink == "" {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	var link string
	for key, value := range db.DB {
		if value == shortLink {
			link = key
			break
		}
	}
	if link == "" {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	w.Header().Set("Location", link)
	w.WriteHeader(http.StatusTemporaryRedirect)
	// w.Write([]byte{})
}
