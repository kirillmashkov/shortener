package config

import (
	"flag"
	"fmt"

	"github.com/caarlos0/env/v6"
)

type serverConfig struct {
	Host     string "env:\"SERVER_ADDRESS\""
	Redirect string "env:\"BASE_URL\""
}

var ServerEnv serverConfig
var ServerArg serverConfig

func init() {
	env.Parse(&ServerEnv)

	fmt.Printf("SERVER_ADDRESS = %s \r\n", ServerEnv.Host)
	fmt.Printf("BASE_URL = %s \r\n", ServerEnv.Redirect)

	flag.StringVar(&ServerArg.Host, "a", "localhost:8080", "server host")
	flag.StringVar(&ServerArg.Redirect, "b", "http://localhost:8080", "server redirect")
}
