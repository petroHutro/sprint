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

func TestHandlerGetUrls(t *testing.T) {
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
		code   int
		answer string
	}
	tests := []struct {
		name string

		request string
		long    []string
		short   []string
		id      []string
		flagB   []bool
		cookie  *http.Cookie
		idUser  string

		want want
	}{
		{
			name: `
GET #1
correct all
got status 200
`,
			request: "/api/user/urls",
			long:    []string{"long1", "long2", "long3"},
			short:   []string{"short1", "short2", "short3"},
			id:      []string{"1", "1", "1"},
			flagB:   []bool{false, false, false},
			cookie:  &http.Cookie{Name: "Authorization", Value: "token"},
			idUser:  "1",
			want: want{
				code:   200,
				answer: `[{"short_url":"http://localhost:8080/short1","original_url":"long1"},{"short_url":"http://localhost:8080/short2","original_url":"long2"},{"short_url":"http://localhost:8080/short3","original_url":"long3"}]`,
			},
		},
		{
			name: `
GET #2
not correct id user in cookie
got status 204
`,
			request: "/api/user/urls",
			long:    []string{"long1", "long2", "long3"},
			short:   []string{"short1", "short2", "short3"},
			id:      []string{"1", "1", "1"},
			flagB:   []bool{false, false, false},
			cookie:  &http.Cookie{Name: "Authorization", Value: "token"},
			idUser:  "2",
			want: want{
				code:   204,
				answer: ``,
			},
		},
		{
			name: `
GET #3
not correct cookie
got status 204
`,
			request: "/api/user/urls",
			long:    []string{"long1", "long2", "long3"},
			short:   []string{"short1", "short2", "short3"},
			id:      []string{"1", "1", "1"},
			flagB:   []bool{false, false, false},
			cookie:  &http.Cookie{Name: "a", Value: "a"},
			idUser:  "1",
			want: want{
				code:   204,
				answer: ``,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
			defer cancel()

			for i, el := range tt.long {
				st.Set(ctx, el, tt.short[i], tt.id[i], tt.flagB[i])
			}

			r := httptest.NewRequest(http.MethodGet, tt.request, nil)
			w := httptest.NewRecorder()
			r.Header.Set("User_id", tt.idUser)
			r.AddCookie(tt.cookie)
			handlers.HandlerGetUrls(w, r, "http://localhost:8080", st)

			rez := w.Result()
			defer rez.Body.Close()

			assert.Equal(t, tt.want.code, rez.StatusCode)
			if rez.StatusCode == 200 {
				assert.Equal(t, tt.want.answer, w.Body.String())
			}
		})
	}
}
