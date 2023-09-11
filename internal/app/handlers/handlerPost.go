package handlers

import (
	"io"
	"net/http"
	"sprint/internal/app/config"
	"sprint/internal/app/storage"
	"sprint/internal/app/utils"
)

func HandlerPost(w http.ResponseWriter, r *http.Request, flag *config.Flags) {
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
	baseAdress := string(flag.BaseURL)
	w.Write([]byte(baseAdress + "/" + storage.GetDB(string(link))))

}
