package handlers_test

import (
	"net/http"
	"net/http/httptest"
	"sprint/cmd/shortener/db"
	"sprint/cmd/shortener/handlers"
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_requestGet(t *testing.T) {
	type want struct {
		code int
		link string
	}
	tests := []struct {
		name      string
		request   string
		shortLink string
		link      string
		want      want
	}{
		{
			name:      "GET#1 Test",
			request:   "/19",
			shortLink: "19",
			link:      "123456789",
			want: want{
				code: 307,
				link: "123456789",
			},
		},
		{
			name:      "GET#2 Test",
			request:   "/",
			shortLink: "",
			link:      "",
			want: want{
				code: 400,
				link: "",
			},
		},
		{
			name:      "GET#3 Test",
			request:   "/1234",
			shortLink: "123",
			link:      "123",
			want: want{
				code: 400,
				link: "",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			db.DB[tt.want.link] = tt.shortLink
			r := httptest.NewRequest(http.MethodGet, tt.request, nil)
			w := httptest.NewRecorder()
			handlers.HandlerGet(w, r)

			rez := w.Result()
			defer rez.Body.Close()
			assert.Equal(t, tt.want.code, rez.StatusCode)
			if rez.StatusCode == 307 {
				assert.Contains(t, rez.Header, "Location")
			}
			assert.Equal(t, tt.want.link, rez.Header.Get("Location"))
		})
	}
}