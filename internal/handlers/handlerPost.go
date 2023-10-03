package handlers

import (
	"io"
	"net/http"
	"sprint/internal/storage"
	"sprint/internal/utils"

	"go.uber.org/zap"
)

func HandlerPost(w http.ResponseWriter, r *http.Request, baseAddress, file string, log *zap.Logger) {
	if r.URL.Path != "/" || utils.ValidContentType(r.Header.Get("Content-Type"), "text/plain") != nil {
		log.Info("Post not Path or not Content-Type")
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	defer r.Body.Close()
	link, err := io.ReadAll(r.Body)
	if err != nil || len(link) == 0 {
		log.Info("Post not body")
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	storage.LongToShort(string(link), file)
	w.Header().Set("Content-Type", "text/plain")
	w.WriteHeader(http.StatusCreated)
	w.Write([]byte(baseAddress + "/" + storage.GetDB(string(link))))
}
