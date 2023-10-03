package handlers

import (
	"io"
	"log"
	"net/http"
	"sprint/internal/storage"
	"sprint/internal/utils"
)

func HandlerPost(w http.ResponseWriter, r *http.Request, baseAddress, file string) {
	if r.URL.Path != "/" || utils.ValidContentType(r.Header.Get("Content-Type"), "text/plain") != nil {
		log.Print("Post not Path or not Content-Type")
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	defer r.Body.Close()
	link, err := io.ReadAll(r.Body)
	if err != nil || len(link) == 0 {
		log.Print("Post not body")
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	storage.LongToShort(string(link), file)
	w.Header().Set("Content-Type", "text/plain")
	w.WriteHeader(http.StatusCreated)
	w.Write([]byte(baseAddress + "/" + storage.GetDB(string(link))))
}
