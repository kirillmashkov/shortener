package main

import (
	"flag"
	"net/http"

	"github.com/kirillmashkov/shortener.git/internal/app"
	"github.com/kirillmashkov/shortener.git/internal/config"
	"github.com/kirillmashkov/shortener.git/internal/httpserver/middleware/logger"
	"github.com/kirillmashkov/shortener.git/internal/httpserver/router"
	"github.com/kirillmashkov/shortener.git/internal/storage"
)

func main() {
	logger.Initialize()
	flag.Parse()
	config.InitServerConf()
	storage.InitStorage()

	err := http.ListenAndServe(app.ServerConf.Host, router.Serv())
	if err != nil {
		panic(err)
	}
}
