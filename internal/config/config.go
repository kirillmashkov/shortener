package config

import (
	"flag"
	"github.com/caarlos0/env/v6"
	"go.uber.org/zap"
)

type ServerConfig struct {
	Host        string "env:\"SERVER_ADDRESS\""
	Redirect    string "env:\"BASE_URL\""
	FileStorage string "env:\"FILE_STORAGE_PATH\""
	Connection	string "env:\"DATABASE_DSN\""
}

var ServerEnv ServerConfig
var ServerArg ServerConfig

func init() {
	flag.StringVar(&ServerArg.Host, "a", "localhost:8080", "server host")
	flag.StringVar(&ServerArg.Redirect, "b", "http://localhost:8080", "server redirect")
	flag.StringVar(&ServerArg.FileStorage, "f", "short_url_storage.txt", "file storage short url")
	flag.StringVar(&ServerArg.Connection, "d", "postgres://postgres:postgres@localhost:5432/postgres", "db connection string")
}

func InitServerConf(conf *ServerConfig, logger *zap.Logger) {
	err := env.Parse(&ServerEnv)
	if err != nil {
		logger.Error("Can't read env variables")
	}

	conf.Redirect = getConfigString(ServerEnv.Redirect, ServerArg.Redirect)
	conf.Host = getConfigString(ServerEnv.Host, ServerArg.Host)
	conf.FileStorage = getConfigString(ServerEnv.FileStorage, ServerArg.FileStorage)
	conf.Connection = getConfigString(ServerEnv.Connection, ServerArg.Connection)

	logger.Info("server config",
		zap.String("host", conf.Host),
		zap.String("redirect", conf.Redirect),
		zap.String("file_storage", conf.FileStorage),
		zap.String("db connection", conf.Connection))
}

func getConfigString(env string, arg string) string {
	if env == "" {
		return arg
	} else {
		return env
	}
}