package handlers_test

import (
	"io"
	"net/http"
	"net/http/httptest"
	"sprint/internal/config"
	"sprint/internal/handlers"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestHandlerPostAPI(t *testing.T) {
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
			name: "POST#1 Test",
			request: request{
				url:         "/api/shorten",
				body:        `{"url": "https://practicum.yandex.ru"}`,
				contentType: "application/json",
			},
			want: want{
				code:        201,
				contentType: "application/json",
			},
		},
		{
			name: "POST#2 Test",
			request: request{
				url:         "/api/shorten",
				body:        "",
				contentType: "application/json",
			},
			want: want{
				code:        400,
				contentType: "",
			},
		},
		{
			name: "POST#3 Test",
			request: request{
				url:         "/api/shorten",
				body:        `{"url": "https://practicum.yandex.ru"}`,
				contentType: "text/plain",
			},
			want: want{
				code:        400,
				contentType: "",
			},
		},
		{
			name: "POST#4 Test",
			request: request{
				url:         "/api",
				body:        `{"url": "https://practicum.yandex.ru"}`,
				contentType: "application/json",
			},
			want: want{
				code:        400,
				contentType: "",
			},
		},
		{
			name: "POST#5 Test",
			request: request{
				url:         "/api/shorten",
				body:        `{"gg": "https://practicum.yandex23.ru"}`,
				contentType: "application/json",
			},
			want: want{
				code:        400,
				contentType: "",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			flags := config.NewFlags()
			body := strings.NewReader(tt.body)
			r := httptest.NewRequest(http.MethodPost, tt.url, body)
			r.Header.Set("Content-Type", tt.contentType)
			w := httptest.NewRecorder()
			handlers.HandlerPostAPI(w, r, string(flags.BaseURL), "")
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
