package handlers

import (
	"database/sql"
	"net/http"
	"sprint/internal/storage"
)

func HandlerPing(w http.ResponseWriter, r *http.Request, db *sql.DB) {
	if err := storage.PingDB(db); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
	}
	w.WriteHeader(http.StatusOK)
}
