package compression

import (
	"bytes"
	"io"
	"net/http"

	"sprint/internal/utils"
	"strings"
)

func GzipMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ow := w
		acceptEncoding := r.Header.Get("Accept-Encoding")
		supportsGzip := strings.Contains(acceptEncoding, "gzip")
		if supportsGzip {
			contentType := r.Header.Get("Content-Type")
			if contentType == "application/json" || contentType == "text/html" {
				cw := newCompressWriter(w)
				ow = cw
				// ow.Header().Set("Content-Type", "application/x")
				// r.Header.Set("Content-Type", "application/x") //мб надо
				defer cw.Close()
			}
		}
		contentEncoding := r.Header.Get("Content-Encoding")
		sendsGzip := strings.Contains(contentEncoding, "gzip")
		if sendsGzip {
			cr, err := newCompressReader(r.Body)
			if err != nil {
				w.WriteHeader(http.StatusInternalServerError)
				return
			}
			data, err := io.ReadAll(cr)
			if err != nil {
				w.WriteHeader(http.StatusInternalServerError)
				return
			}
			r.Body = io.NopCloser(bytes.NewReader(data))
			// fmt.Println(string(data))
			if utils.IsJSON(data) {
				r.Header.Set("Content-Type", "application/json")
				// fmt.Println("application/json")
			} else if utils.IsText(data) {
				// fmt.Println("text/plain")
				r.Header.Set("Content-Type", "text/plain")
			} else {
				w.WriteHeader(http.StatusInternalServerError)
				return
			}
			// contentType := http.DetectContentType(data)
			// fmt.Println(contentType)
			// if utils.ValidContentType(contentType, "application/json") != nil && utils.ValidContentType(contentType, "text/plain") != nil {
			// 	w.WriteHeader(http.StatusInternalServerError)
			// 	return
			// }
			// r.Header.Set("Content-Type", contentType)
			defer cr.Close()
		}
		next.ServeHTTP(ow, r)
	})
}
