package handlers_test

import (
	"io"
	"net/http"
	"net/http/httptest"
	"sprint/internal/config"
	"sprint/internal/handlers"
	"sprint/internal/logger"
	"sprint/internal/storage"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func Test_requestPost(t *testing.T) {
	flags := config.NewFlags()
	if err := logger.InitLogger(flags.Logger); err != nil {
		logger.Panic(err.Error())
	}
	defer logger.Shutdown()

	confData := config.Storage{
		FileStoragePath: "",
		DatabaseDSN:     "",
	}
	st, _ := storage.InitStorage(&confData)

	type request struct {
		url         string
		body        string
		contentType string
	}
	type want struct {
		code        int
		contentType string
	}
	tests := []struct {
		name string
		request
		want want
	}{
		{
			name: `
POST / #1 
correct url, correct body, correct contentType
got status 201
`,
			request: request{
				url:         "/",
				body:        "https://www.google.com/",
				contentType: "text/plain; charset=utf-8",
			},
			want: want{
				code:        201,
				contentType: "text/plain",
			},
		},
		{
			name: `
POST / #2 
correct url, empty body, correct contentType
got status 400
`,
			request: request{
				url:         "/",
				body:        "",
				contentType: "text/plain; charset=utf-8",
			},
			want: want{
				code:        400,
				contentType: "",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			body := strings.NewReader(tt.body)
			r := httptest.NewRequest(http.MethodPost, tt.url, body)
			r.Header.Set("Content-Type", tt.contentType)
			w := httptest.NewRecorder()

			handlers.HandlerPost(w, r, flags.BaseURL, flags.FileStoragePath, st)
			rez := w.Result()
			defer rez.Body.Close()

			assert.Equal(t, tt.want.code, rez.StatusCode)
			assert.Equal(t, tt.want.contentType, rez.Header.Get("Content-Type"))

			rezBody, err := io.ReadAll(rez.Body)
			require.NoError(t, err)
			err = rez.Body.Close()
			require.NoError(t, err)

			if rez.StatusCode == 201 {
				assert.NotEmpty(t, string(rezBody))
			} else if rez.StatusCode == 400 {
				assert.Empty(t, string(rezBody))
			}
		})
	}
}
