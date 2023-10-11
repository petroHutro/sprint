package handlers

import (
	"net/http"
	"sprint/internal/logger"
	"sprint/internal/storage"
)

func HandlerGet(w http.ResponseWriter, r *http.Request) {
	shortLink := r.URL.String()[1:]
	if shortLink == "" {
		logger.Log.Error("shortLink is emty")
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	link, err := storage.ShortToLong(shortLink)
	if err != nil {
		logger.Log.Error("cannot convert short to long :%v", err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	w.Header().Set("Location", link)
	w.WriteHeader(http.StatusTemporaryRedirect)
}
