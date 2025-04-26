package main

import (
	"flag"
	"net/http"

	"github.com/kirillmashkov/shortener.git/internal/app"
	"github.com/kirillmashkov/shortener.git/internal/logger"
	"github.com/kirillmashkov/shortener.git/internal/httpserver/router"
)

func main() {
	err := logger.Initialize()
	if err != nil {
		panic(err)
	}

	flag.Parse()
	err = app.Initialize()
	if err != nil {
		panic(err)
	}

	err = http.ListenAndServe(app.ServerConf.Host, router.Serv())
	if err != nil {
		panic(err)
	}
}
