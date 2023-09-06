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

	defer r.Body.Close()
	link, err := io.ReadAll(r.Body)
	if err != nil || len(link) == 0 {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	//возможно пересмотреть порядок
	if _, ok := db.DB[string(link)]; !ok {
		shortLink := utils.LinkShortening()
		db.DB[string(link)] = shortLink
	}

	w.Header().Set("Content-Type", "text/plain")
	w.WriteHeader(http.StatusCreated)
	w.Write([]byte(db.DB[string(link)]))
}

func HandlerGet(w http.ResponseWriter, r *http.Request) {
	shortLink := r.URL.String()[1:]
	if shortLink == "" {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	var link string
	for key, value := range db.DB {
		if value == shortLink {
			link = key
			break
		}
	}
	if link == "" {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	w.Header().Set("Location", link)
	w.WriteHeader(http.StatusTemporaryRedirect)
	// w.Write([]byte{})
}
