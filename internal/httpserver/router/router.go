package router

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/kirillmashkov/shortener.git/internal/httpserver/handler"
	"github.com/kirillmashkov/shortener.git/internal/httpserver/middleware/logger"
)

func Serv() http.Handler {
	logger.Initialize()

	r := chi.NewRouter()
	r.Use(logger.Logger)
	r.Post("/", handler.PostHandler)
	r.Get("/{id}", handler.GetHandler)
	r.Post("/shorten", handler.PostGenerateShortURL)
	return r
}