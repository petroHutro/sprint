package handlers

import (
	"io"
	"net/http"
	"sprint/internal/logger"
	"sprint/internal/storage"
)

func HandlerPost(w http.ResponseWriter, r *http.Request, baseAddress, file string, db *storage.StorageBase) {
	link, err := io.ReadAll(r.Body)
	if err != nil || len(link) == 0 {
		logger.Log.Error("Post not body %v", err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	defer r.Body.Close()
	err = db.LongToShort(r.Context(), string(link), file)
	statusCode := http.StatusCreated
	if err != nil {
		if repErr, ok := err.(*storage.RepError); ok && repErr.Repetition {
			statusCode = http.StatusConflict
			logger.Log.Error("long already db :%v", err)
		} else {
			logger.Log.Error("cannot convert long to short :%v", err)
			w.WriteHeader(http.StatusBadRequest)
			return
		}
	}
	w.Header().Set("Content-Type", "text/plain")
	w.WriteHeader(statusCode)
	w.Write([]byte(baseAddress + "/" + db.GetShort(r.Context(), string(link))))
}
