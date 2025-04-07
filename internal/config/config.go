package config

import (
	"flag"
	"fmt"
	"github.com/caarlos0/env/v6"
	"github.com/kirillmashkov/shortener.git/internal/app"
)



func init() {
	env.Parse(&app.ServerEnv)

	fmt.Printf("SERVER_ADDRESS = %s \r\n", app.ServerEnv.Host)
	fmt.Printf("BASE_URL = %s \r\n", app.ServerEnv.Redirect)

	flag.StringVar(&app.ServerArg.Host, "a", "localhost:8080", "server host")
	flag.StringVar(&app.ServerArg.Redirect, "b", "http://localhost:8080", "server redirect")
}

func InitServerConf() {
	app.ServerConf.Host = getConfigString(app.ServerEnv.Host, app.ServerArg.Host)
	app.ServerConf.Redirect = getConfigString(app.ServerEnv.Redirect, app.ServerArg.Redirect)

	fmt.Printf("host = %s, redirect = %s \r\n", app.ServerConf.Host, app.ServerConf.Redirect)
}

func getConfigString(env string, arg string) string {
	if env == "" {
		return arg
	} else {
		return env
	}
}

