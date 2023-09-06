package handlers

import (
	"io"
	"net/http"
	"sprint/cmd/shortener/db"
	"sprint/cmd/shortener/utils"
)

func HandlerPost(w http.ResponseWriter, r *http.Request) {
	//пределать это условие с учетом всех text/plain
	if r.URL.Path != "/" || r.Header.Get("Content-Type") != "text/plain; charset=utf-8" {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	link, err := io.ReadAll(r.Body)
	if err != nil || len(link) == 0 {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	if _, ok := db.DB[string(link)]; !ok {
		shortLink := utils.LinkShortening()
		db.DB[string(link)] = shortLink
	}

	w.Header().Set("Content-Type", "text/plain")
	w.WriteHeader(http.StatusCreated)
	w.Write([]byte(db.DB[string(link)]))
}
