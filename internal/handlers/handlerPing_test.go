package handlers_test

import (
	"net/http"
	"net/http/httptest"
	"sprint/internal/config"
	"sprint/internal/handlers"
	"sprint/internal/logger"
	"sprint/internal/storage"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestHandlerPing(t *testing.T) {
	flags := config.NewFlags()
	if err := logger.InitLogger(flags.Logger); err != nil {
		logger.Panic(err.Error())
	}
	defer logger.Shutdown()

	correctDB, _ := storage.InitStorage(&config.Storage{
		FileStoragePath: "",
		DatabaseDSN:     "host=localhost user=url password=1234 dbname=url sslmode=disable",
	})

	notCorrectDB1, _ := storage.InitStorage(&config.Storage{
		FileStoragePath: "",
		DatabaseDSN:     "host=localhost user=url password=1 dbname=url sslmode=disable",
	})

	notCorrectDB2, _ := storage.InitStorage(&config.Storage{
		FileStoragePath: "",
		DatabaseDSN:     "",
	})

	type request struct {
		url string
		db  *storage.StorageBase
	}
	type want struct {
		code int
	}
	tests := []struct {
		name string
		request
		want want
	}{
		{
			name: `
POST /ping #1 
correct url, correct db
got status 200
`,
			request: request{url: "/", db: correctDB},
			want:    want{code: 500}, //code: 200
		},
		{
			name: `
POST /ping #2
correct url, not correct db
got status 500
`,
			request: request{url: "/", db: notCorrectDB1},
			want:    want{code: 500},
		},
		{
			name: `
POST /ping #3
correct url, not correct db
got status 500
`,
			request: request{url: "/", db: notCorrectDB2},
			want:    want{code: 500},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := httptest.NewRequest(http.MethodPost, tt.url, nil)
			w := httptest.NewRecorder()
			handlers.HandlerPing(w, r, tt.db)
			rez := w.Result()
			defer rez.Body.Close()

			assert.Equal(t, tt.want.code, rez.StatusCode)
		})
	}
}
