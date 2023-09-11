package handlers

import (
	"net/http"
	"sprint/internal/app/storage"
)

func HandlerGet(w http.ResponseWriter, r *http.Request) {
	shortLink := r.URL.String()[1:]
	if shortLink == "" {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	link, err := storage.ShortToLong(shortLink)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	w.Header().Set("Location", link)
	w.WriteHeader(http.StatusTemporaryRedirect)
}
