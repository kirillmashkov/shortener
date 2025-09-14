package httpserver

import (
	"context"
	"net/http"

	"github.com/kirillmashkov/shortener.git/internal/httpserver/router"
	"github.com/kirillmashkov/shortener.git/internal/server"
	"golang.org/x/crypto/acme/autocert"
)

type HTTPS struct {
	server *http.Server
}

// Run - запуск https сервера
func (s *HTTPS) Run() error {
	return s.server.ListenAndServe()
}

// Shutdown - остановка https сервера
func (s *HTTPS) Shutdown() error {
	return s.server.Shutdown(context.Background())
}

func NewHTTPS(addr string) server.Server {
	httpServer := &http.Server{
		Addr:              addr,
		Handler:           router.Serv(),
	}

	manager := &autocert.Manager{
		Cache:      autocert.DirCache("cache-dir"),
		Prompt:     autocert.AcceptTOS,
		HostPolicy: autocert.HostWhitelist("mysite.ru", "www.mysite.ru"),
	}

	httpServer.TLSConfig = manager.TLSConfig()

	server := &HTTPS{
		server: httpServer,
	}
	
	return server
}