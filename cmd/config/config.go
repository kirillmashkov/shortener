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
	flag.StringVar(&Server.Redirect, "b", "http://localhost:8080", "server redirect")
}