package main

import (
	"flag"
	"net/http"

	"github.com/kirillmashkov/shortener.git/internal/config"
	"github.com/kirillmashkov/shortener.git/internal/service"
)

func main() {
	flag.Parse()
	config.InitServerConf()

	err := http.ListenAndServe(config.ServerConf.Host, service.Serv())
	if err != nil {
		panic(err)
	}
}
