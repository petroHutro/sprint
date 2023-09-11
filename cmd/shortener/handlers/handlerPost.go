package handlers

import (
	"io"
	"net/http"
	"sprint/cmd/shortener/config"
	"sprint/cmd/shortener/db"
	"sprint/cmd/shortener/utils"
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
	db.LongToShort(string(link))
	w.Header().Set("Content-Type", "text/plain")
	w.WriteHeader(http.StatusCreated)
	baseAdress := string(flag.BaseURL)
	w.Write([]byte(baseAdress + db.GetDB(string(link))))

}
