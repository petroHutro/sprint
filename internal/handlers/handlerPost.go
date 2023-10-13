package handlers

import (
	"io"
	"net/http"
	"sprint/internal/logger"
	"sprint/internal/storage"
)

func HandlerPost(w http.ResponseWriter, r *http.Request, baseAddress, file string, db *storage.StorageBase) {
	// if r.URL.Path != "/" || utils.ValidContentType(r.Header.Get("Content-Type"), "text/plain") != nil {
	// 	log.Print("Post not Path or not Content-Type")
	// 	w.WriteHeader(http.StatusBadRequest)
	// 	return
	// }
	link, err := io.ReadAll(r.Body)
	if err != nil || len(link) == 0 {
		logger.Log.Error("Post not body %v", err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	defer r.Body.Close()
	if err := db.LongToShort(r.Context(), string(link), file); err != nil {
		logger.Log.Error("cannot convert long to short :%v", err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	w.Header().Set("Content-Type", "text/plain")
	w.WriteHeader(http.StatusCreated)
	w.Write([]byte(baseAddress + "/" + db.GetShort(r.Context(), string(link))))
}
