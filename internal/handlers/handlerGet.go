package handlers

import (
	"errors"
	"net/http"
	"sprint/internal/logger"
	"sprint/internal/storage"
)

func HandlerGet(w http.ResponseWriter, r *http.Request, db *storage.StorageBase) {
	shortLink := r.URL.String()[1:]
	if shortLink == "" {
		logger.Error("shortLink is emty")
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	link, err := db.ShortToLong(r.Context(), shortLink)
	if err != nil {
		if err == errors.New("url delete") {
			logger.Error("link alrege delete:%v", err)
			w.WriteHeader(http.StatusGone)
			return
		}
		logger.Error("cannot convert short to long :%v", err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	w.Header().Set("Location", link)
	w.WriteHeader(http.StatusTemporaryRedirect)
}
