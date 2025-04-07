package router

import (
	"github.com/go-chi/chi/v5"
	"net/http"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/kirillmashkov/shortener.git/internal/httpserver/handler"
 )

func Serv() http.Handler {
	r := chi.NewRouter()
	r.Use(middleware.Logger)
	r.Post("/", handler.PostHandler)
	r.Get("/{id}", handler.GetHandler)
	return r
}