package handler

import (
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestPostHandler(t *testing.T) {
	type want struct {
		code        int
		contentType string
	}

	tests := []struct {
		name        string
		requestLink io.Reader
		want        want
	}{
		{
			name:        "test successful create short link",
			requestLink: strings.NewReader("https://www.lenta.ru"),
			want: want{
				code:        201,
				contentType: "text/plain",
			},
		},
		{
			name:        "test empty link",
			requestLink: nil,
			want: want{
				code:        400,
				contentType: "text/plain",
			},
		},
		{
			name:        "test bad link",
			requestLink: strings.NewReader("blablabla"),
			want: want{
				code:        400,
				contentType: "text/plain",
			},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {

			request := httptest.NewRequest(http.MethodPost, "/", test.requestLink)
			// создаём новый Recorder
			w := httptest.NewRecorder()
			PostHandler(w, request)

			res := w.Result()
			// проверяем код ответа
			assert.Equal(t, test.want.code, res.StatusCode)

			// получаем и проверяем тело запроса
			_, err := io.ReadAll(res.Body)
			defer res.Body.Close()
			require.NoError(t, err)
		})
	}
}

func TestGetHandler(t *testing.T) {
	type want struct {
		code        int
		contentType string
	}

	tests := []struct {
		name string
		want want
	}{
		{
			name: "test status 400, unknown key",
			want: want{
				code:        400,
				contentType: "text/plain",
			},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			request := httptest.NewRequest(http.MethodGet, "/aaaaa", nil)
			// создаём новый Recorder
			w := httptest.NewRecorder()
			GetHandler(w, request)

			res := w.Result()
			// проверяем код ответа
			assert.Equal(t, test.want.code, res.StatusCode)

			// получаем и проверяем тело запроса
			_, err := io.ReadAll(res.Body)
			defer res.Body.Close()
			require.NoError(t, err)
		})
	}
}
