package handlers

import (
	"io"
	"net/http"
	"sprint/cmd/shortener/db"
	"sprint/cmd/shortener/utils"
)

func HandlerPost(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" || r.Header.Get("Content-Type") != "text/plain" {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	link, err := io.ReadAll(r.Body)
	if err != nil || len(link) == 0 {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	if _, ok := db.Db[string(link)]; !ok {
		shortLink := utils.LinkShortening()
		db.Db[string(link)] = shortLink
	}

	w.Header().Set("Content-Type", "text/plain")
	w.WriteHeader(http.StatusCreated)
	w.Write([]byte(db.Db[string(link)]))
}
