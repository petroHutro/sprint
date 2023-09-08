package handlers

import (
	"net/http"
	"sprint/cmd/shortener/db"
)

func HandlerGet(w http.ResponseWriter, r *http.Request) {
	shortLink := r.URL.String()[1:]
	if shortLink == "" {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	link, err := db.ShortToLong(shortLink)
	if err {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	w.Header().Set("Location", link)
	w.WriteHeader(http.StatusTemporaryRedirect)
}
