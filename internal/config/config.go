package config

import (
	"flag"
	"github.com/caarlos0/env/v6"
	"github.com/kirillmashkov/shortener.git/internal/app"
	"go.uber.org/zap"
)

func init() {
	env.Parse(&app.ServerEnv)

	flag.StringVar(&app.ServerArg.Host, "a", "localhost:8080", "server host")
	flag.StringVar(&app.ServerArg.Redirect, "b", "http://localhost:8080", "server redirect")
	flag.StringVar(&app.ServerArg.FileStorage, "f", "short_url_storage.txt", "file storage short url")
}

func InitServerConf() {
	app.ServerConf.Host = getConfigString(app.ServerEnv.Host, app.ServerArg.Host)
	app.ServerConf.Redirect = getConfigString(app.ServerEnv.Redirect, app.ServerArg.Redirect)
	app.ServerConf.FileStorage = getConfigString(app.ServerEnv.FileStorage, app.ServerArg.FileStorage)

	app.Log.Info("server config",
		zap.String("host", app.ServerConf.Host),
		zap.String("redirect", app.ServerConf.Redirect),
		zap.String("file_storage", app.ServerConf.FileStorage))
}

func getConfigString(env string, arg string) string {
	if env == "" {
		return arg
	} else {
		return env
	}
}

