package config

import (
	"flag"

	"github.com/caarlos0/env/v6"
	"go.uber.org/zap"
)

// ServerConfig - тип для хранения конфигурации приложения
type ServerConfig struct {
	Host        string "env:\"SERVER_ADDRESS\""
	Redirect    string "env:\"BASE_URL\""
	FileStorage string "env:\"FILE_STORAGE_PATH\""
	Connection  string "env:\"DATABASE_DSN\""
	EnableHTTPS bool   "env:\"ENABLE_HTTPS\""
}

// ServerEnv - хранение значений, полученных из переменных среды
var ServerEnv ServerConfig

// ServerArg - хранение значений, полученных из командной строки
var ServerArg ServerConfig

func init() {
	flag.StringVar(&ServerArg.Host, "a", "localhost:8080", "server host")
	flag.StringVar(&ServerArg.Redirect, "b", "http://localhost:8080", "server redirect")
	flag.StringVar(&ServerArg.FileStorage, "f", "short_url_storage.txt", "file storage short url")
	flag.StringVar(&ServerArg.Connection, "d", "postgres://postgres:postgres@localhost:5432/postgres?sslmode=disable", "db connection string")
	flag.BoolVar(&ServerArg.EnableHTTPS, "s", false, "run server with https")
}

// InitServerConf - определение итоговой конфигурации приложения
func InitServerConf(conf *ServerConfig, logger *zap.Logger) {
	err := env.Parse(&ServerEnv)
	if err != nil {
		logger.Error("Can't read env variables")
	}

	conf.Redirect = getConfigString(ServerEnv.Redirect, ServerArg.Redirect)
	conf.Host = getConfigString(ServerEnv.Host, ServerArg.Host)
	conf.FileStorage = getConfigString(ServerEnv.FileStorage, ServerArg.FileStorage)
	conf.Connection = getConfigString(ServerEnv.Connection, ServerArg.Connection)
	conf.EnableHTTPS = getConfigBool(ServerEnv.EnableHTTPS, ServerArg.EnableHTTPS)

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
func getConfigBool(env bool, arg bool) bool {
	if !env {
		return arg
	} else {
		return env
	}
}
