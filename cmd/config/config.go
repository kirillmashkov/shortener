package config

import (
	"flag"
)

type serverConfig struct {
	Host string
	Redirect string
}

var Server serverConfig

func Init() {
	flag.StringVar(&Server.Host, "a", "localhost:8080", "server host")
	flag.StringVar(&Server.Redirect, "b", "localhost:8080", "server redirect")
}