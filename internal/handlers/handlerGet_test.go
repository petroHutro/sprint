package handlers_test

import (
	"context"
	"net/http"
	"net/http/httptest"
	"sprint/internal/config"
	"sprint/internal/handlers"
	"sprint/internal/logger"
	"sprint/internal/storage"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func Test_requestGet(t *testing.T) {
	log := config.Logger{
		FilePath:  "file.log",
		FileFlag:  false,
		MultiFlag: false,
	}
	if err := logger.InitLogger(log); err != nil {
		logger.Panic(err.Error())
	}
	defer logger.Shutdown()

	confData := config.Storage{
		FileStoragePath: "",
		DatabaseDSN:     "",
	}
	st, _ := storage.InitStorage(&confData)

	type want struct {
		code int
		link string
	}
	tests := []struct {
		name      string
		request   string
		shortLink string
		link      string
		flagB     bool
		want      want
	}{
		{
			name: `
GET #1
correct url, correct longLink, correct shortLink
got status 307
`,
			request:   "/19",
			shortLink: "19",
			link:      "123456789",
			flagB:     false,
			want: want{
				code: 307,
				link: "123456789",
			},
		},
		{
			name: `
GET #2
emty url, emty longLink, emty shortLink
got status 400
`,
			request:   "/",
			shortLink: "",
			link:      "",
			flagB:     false,
			want: want{
				code: 400,
				link: "",
			},
		},
		{
			name: `
GET #3
correct url, not correct longLink, not correct shortLink
got status 400
`,
			request:   "/1234",
			shortLink: "123",
			link:      "123",
			flagB:     false,
			want: want{
				code: 400,
				link: "",
			},
		},
		{
			name: `
GET #4
delete url,
got status 400
`,
			request:   "/199",
			shortLink: "199",
			link:      "1234567890",
			flagB:     true,
			want: want{
				code: 410,
				link: "",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
			defer cancel()
			st.Set(ctx, tt.link, tt.shortLink, "1", tt.flagB)
			r := httptest.NewRequest(http.MethodGet, tt.request, nil)
			w := httptest.NewRecorder()

			handlers.HandlerGet(w, r, st)
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
