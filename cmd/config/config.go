package config

import (
	"flag"
	"github.com/caarlos0/env/v6"
)

type serverConfig struct {
	Host string "env:\"BASE_ADDRESS\""
	Redirect string "env:\"BASE_URL\""
}

var Server serverConfig

func init() {
	env.Parse(&Server)

	if Server.Host == "" {
		flag.StringVar(&Server.Host, "a", "localhost:8080", "server host")
	}

	if Server.Redirect == "" {
		flag.StringVar(&Server.Redirect, "b", "http://localhost:8080", "server redirect")
	}
}