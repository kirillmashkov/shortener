package router

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/kirillmashkov/shortener.git/internal/httpserver/handler"
	"github.com/kirillmashkov/shortener.git/internal/httpserver/middleware/logger"
	"github.com/kirillmashkov/shortener.git/internal/httpserver/middleware/compress"
)

func Serv() http.Handler {
	r := chi.NewRouter()
	r.Use(logger.Logger)
	r.Use(compress.Compress)
	r.Post("/", handler.PostHandler)
	r.Get("/{id}", handler.GetHandler)
	r.Post("/api/shorten", handler.PostGenerateShortURL)
	return r
}