package main

import (
	"flag"
	"net/http"

	"github.com/kirillmashkov/shortener.git/internal/config"
	"github.com/kirillmashkov/shortener.git/internal/app"
	"github.com/kirillmashkov/shortener.git/internal/httpserver/router"
)

func main() {
	flag.Parse()
	config.InitServerConf()

	err := http.ListenAndServe(app.ServerConf.Host, router.Serv())
	if err != nil {
		panic(err)
	}
}
