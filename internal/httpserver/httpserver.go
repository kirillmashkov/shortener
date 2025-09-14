package httpserver

import (
	"context"
	"net/http"

	"github.com/kirillmashkov/shortener.git/internal/httpserver/router"
	"github.com/kirillmashkov/shortener.git/internal/server"
)

// HTTP - хранение HTTP сервера
type HTTP struct {
	server *http.Server
}

// Run - запуск http сервера
func (s *HTTP) Run() error {
	return s.server.ListenAndServe()
}

// Shutdown - остановка http сервера
func (s *HTTP) Shutdown() error {
	return s.server.Shutdown(context.Background())
}

func NewHTTP(addr string) server.Server {
	httpServer := &http.Server{
		Addr:    addr,
		Handler: router.Serv(),
	}

	server := &HTTP{
		server: httpServer,
	}

	return server
}
