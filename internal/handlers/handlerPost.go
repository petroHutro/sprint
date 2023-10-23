package handlers

import (
	"errors"
	"io"
	"net/http"
	"sprint/internal/logger"
	"sprint/internal/storage"
	"strconv"
)

func HandlerPost(w http.ResponseWriter, r *http.Request, baseAddress, file string, db *storage.StorageBase) {
	link, err := io.ReadAll(r.Body)
	if err != nil || len(link) == 0 {
		logger.Error("Post not body: %v", err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	defer r.Body.Close()
	id, err := strconv.Atoi(r.Header.Get("User_id"))
	if err != nil {
		logger.Error("bad user id: %v", err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	err = db.LongToShort(r.Context(), string(link), file, id)
	statusCode := http.StatusCreated
	if err != nil {
		var repErr *storage.RepError
		if errors.As(err, &repErr) && repErr.Repetition {
			statusCode = http.StatusConflict
			logger.Error("long already db :%v", err)
		} else {
			logger.Error("cannot convert long to short :%v", err)
			w.WriteHeader(http.StatusBadRequest)
			return
		}
	}
	w.Header().Set("Content-Type", "text/plain")
	w.WriteHeader(statusCode)
	w.Write([]byte(baseAddress + "/" + db.GetShort(r.Context(), string(link))))
}
