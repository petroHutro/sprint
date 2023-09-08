package handlers

import (
	"fmt"
	"io"
	"net/http"
	"sprint/cmd/shortener/db"
	"sprint/cmd/shortener/utils"
)

func HandlerPost(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" || !utils.ValidContentType(r.Header.Get("Content-Type")) {
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
	domain := fmt.Sprintf("http://%s/", r.Host)
	// w.Write([]byte(domain + db.DB[string(link)]))
	w.Write([]byte(domain + db.GetDB(string(link))))

}
