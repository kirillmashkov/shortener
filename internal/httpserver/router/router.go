package router

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/kirillmashkov/shortener.git/internal/httpserver/handler"
	"github.com/kirillmashkov/shortener.git/internal/httpserver/middleware/compress"
	"github.com/kirillmashkov/shortener.git/internal/httpserver/middleware/logger"
	"github.com/kirillmashkov/shortener.git/internal/httpserver/middleware/security"
)

func Serv() http.Handler {
	r := chi.NewRouter()
	r.Use(logger.Logger)
	r.Use(compress.Compress)
	r.Use(security.Auth)
	r.Post("/", handler.PostHandler)
	r.Get("/{id}", handler.GetHandler)
	r.Get("/api/user/urls", handler.GetAllURL)
	r.Post("/api/shorten", handler.PostGenerateShortURL)
	r.Post("/api/shorten/batch", handler.PostGenerateShortURLBatch)
	r.Delete("/api/user/urls", handler.DeleteURLBatch)
	r.Get("/ping", handler.Ping)

	return r
}