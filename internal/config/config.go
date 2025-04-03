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

var serverEnv serverConfig
var serverArg serverConfig
var ServerConf serverConfig

func init() {
	env.Parse(&serverEnv)

	fmt.Printf("SERVER_ADDRESS = %s \r\n", serverEnv.Host)
	fmt.Printf("BASE_URL = %s \r\n", serverEnv.Redirect)

	flag.StringVar(&serverArg.Host, "a", "localhost:8080", "server host")
	flag.StringVar(&serverArg.Redirect, "b", "http://localhost:8080", "server redirect")
}

func InitServerConf() {
	ServerConf.Host = getConfigString(serverEnv.Host, serverArg.Host)
	ServerConf.Redirect = getConfigString(serverEnv.Redirect, serverArg.Redirect)

	fmt.Printf("host = %s, redirect = %s \r\n", ServerConf.Host, ServerConf.Redirect)
}

func getConfigString(env string, arg string) string {
	if env == "" {
		return arg
	} else {
		return env
	}
}

