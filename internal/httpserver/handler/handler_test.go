package handler

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/go-chi/chi/v5"
	"github.com/kirillmashkov/shortener.git/internal/app"
	"github.com/kirillmashkov/shortener.git/internal/httpserver/middleware/compress"
	"github.com/kirillmashkov/shortener.git/internal/httpserver/middleware/security"
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
	}

	err := app.Initialize()
	if err != nil {
		log.SetPrefix("ERROR")
		log.Println("Can't initialize app")
		return
	}

	log.SetFlags(log.Ldate | log.Ltime | log.Lshortfile)
	
	r := chi.NewRouter()
	r.Use(security.Auth)
	r.Post("/", PostHandler)

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {

			request := httptest.NewRequest(http.MethodPost, "/", test.requestLink)
			// создаём новый Recorder
			w := httptest.NewRecorder()
			// PostHandler(w, request)
			r.ServeHTTP(w, request)

			res := w.Result()
			// проверяем код ответа
			assert.Equal(t, test.want.code, res.StatusCode)

			// получаем и проверяем тело запроса
			_, err := io.ReadAll(res.Body)
			defer func() {
				if errClose := res.Body.Close(); errClose != nil {
					log.SetPrefix("ERROR")
					log.Println("Can't close io stream in test")
				}
			}()

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
			defer func() {
				if errClose := res.Body.Close(); errClose != nil {
					fmt.Println("Can't close")
				}
			}()
			require.NoError(t, err)
		})
	}
}

func TestPostGenerateShortURL(t *testing.T) {
	tests := []struct {
		name string
		body string
		expectedCode int
		compress bool
	}{
		{
			name: "test successful create short link",
			body: `{"url": "http://www.lenta.ru"}`,
			expectedCode: 201,
			compress: false,
		},
		{
			name: "test successful create short link",
			body: `{"url": "http://www.lenta.ru"}`,
			expectedCode: 201,
			compress: true,
		},
	}

	r := chi.NewRouter()
	r.Use(compress.Compress)
	r.Use(security.Auth)
	r.Post("/api/shorten", PostGenerateShortURL)

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			request := httptest.NewRequest(http.MethodPost, "/api/shorten", strings.NewReader(test.body))
			if test.compress {
				request.Header.Set("Accept-Encoding", "gzip")
			}

			// создаём новый Recorder
			w := httptest.NewRecorder()

			r.ServeHTTP(w, request)
			res := w.Result()

			// проверяем код ответа
			assert.Equal(t, test.expectedCode, res.StatusCode)
			if test.compress {
				ce := res.Header.Get("Content-Encoding")
				assert.Equal(t, true, strings.Contains(ce, "gzip"))	
			} else {
				ce := res.Header.Get("Content-Encoding")
				assert.Equal(t, true, !strings.Contains(ce, "gzip"))	
			}
			

			// получаем и проверяем тело запроса
			_, err := io.ReadAll(res.Body)
			defer func() {
				if errClose := res.Body.Close(); errClose != nil {
					fmt.Println("Can't close")
				}
			}()
			require.NoError(t, err)
		})
	}	
}
