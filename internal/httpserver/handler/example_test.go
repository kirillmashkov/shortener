package handler

import (
	"log"
	"net/http"
	"net/http/httptest"
	"strings"

	"github.com/go-chi/chi/v5"
	"github.com/kirillmashkov/shortener.git/internal/app"
	"github.com/kirillmashkov/shortener.git/internal/httpserver/middleware/security"
)

func ExamplePostHandler() {
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

	request := httptest.NewRequest(http.MethodPost, "/", strings.NewReader("https://www.lenta.ru"))
	w := httptest.NewRecorder()
	r.ServeHTTP(w, request)
}
