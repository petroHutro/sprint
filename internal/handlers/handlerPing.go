package handlers

import (
	"net/http"
	"sprint/internal/logger"
	"sprint/internal/storage"
)

func HandlerPing(w http.ResponseWriter, r *http.Request, db *storage.StorageBase) {
	if err := db.PingDB(r.Context()); err != nil {
		logger.Error("cannot ping %v", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
}
