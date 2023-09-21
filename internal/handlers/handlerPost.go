package handlers

import (
	"io"
	"net/http"
	"sprint/internal/storage"
	"sprint/internal/utils"
)

func HandlerPost(w http.ResponseWriter, r *http.Request, baseAddress string) {
	if r.URL.Path != "/" || utils.ValidContentType(r.Header.Get("Content-Type")) != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	defer r.Body.Close()
	link, err := io.ReadAll(r.Body)
	if err != nil || len(link) == 0 {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	storage.LongToShort(string(link))
	w.Header().Set("Content-Type", "text/plain")
	w.WriteHeader(http.StatusCreated)
	w.Write([]byte(baseAddress + "/" + storage.GetDB(string(link))))

}
