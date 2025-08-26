package config

import (
	"encoding/json"
	"errors"
	"flag"
	"os"

	"github.com/caarlos0/env/v6"
	"go.uber.org/zap"
)

// ConfigFromFile - тип для хранения конфигурации из файла
type ConfigFromFile struct {
	ServerAddress   string `json:"server_address"`
	BaseURL         string `json:"base_url"`
	FileStoragePath string `json:"file_storage_path"`
	DatabaseDSN     string `json:"database_dsn"`
	EnableHTTPS     bool   `json:"enable_https"`
}

// ServerConfig - тип для хранения конфигурации приложения
type ServerConfig struct {
	Host        string "env:\"SERVER_ADDRESS\""
	Redirect    string "env:\"BASE_URL\""
	FileStorage string "env:\"FILE_STORAGE_PATH\""
	Connection  string "env:\"DATABASE_DSN\""
	EnableHTTPS bool   "env:\"ENABLE_HTTPS\""
	ConfigPath  string "env:\"CONFIG\""
}

const filenameConfigServer = "config/configserver.json"

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
	flag.StringVar(&ServerArg.ConfigPath, "c", "", "config path")
}

// InitServerConf - определение итоговой конфигурации приложения
func InitServerConf(conf *ServerConfig, logger *zap.Logger) {
	err := env.Parse(&ServerEnv)
	if err != nil {
		logger.Error("Can't read env variables")
	}

	conf.ConfigPath = getConfigString(ServerEnv.ConfigPath, ServerArg.ConfigPath, filenameConfigServer)

	var configFromFile *ConfigFromFile
	configFromFile, err = parseConfigFile(conf.ConfigPath)
	if err != nil {
		configFromFile = &ConfigFromFile {
			ServerAddress: "",
			BaseURL: "",
			FileStoragePath: "",
			DatabaseDSN: "",
			EnableHTTPS: false,
		}
	}

	conf.Redirect = getConfigString(ServerEnv.Redirect, ServerArg.Redirect, configFromFile.BaseURL)
	conf.Host = getConfigString(ServerEnv.Host, ServerArg.Host, configFromFile.ServerAddress)
	conf.FileStorage = getConfigString(ServerEnv.FileStorage, ServerArg.FileStorage, configFromFile.FileStoragePath)
	conf.Connection = getConfigString(ServerEnv.Connection, ServerArg.Connection, configFromFile.DatabaseDSN)
	conf.EnableHTTPS = getConfigBool(ServerEnv.EnableHTTPS, ServerArg.EnableHTTPS, configFromFile.EnableHTTPS)

	logger.Info("server config",
		zap.String("host", conf.Host),
		zap.String("redirect", conf.Redirect),
		zap.String("file_storage", conf.FileStorage),
		zap.String("db connection", conf.Connection))
}

func getConfigString(env string, arg string, fromFile string) string {
	if env == "" {
		if arg == "" {
			return fromFile
		} else {
			return arg
		}
	} else {
		return env
	}
}

func getConfigBool(env bool, arg bool, fromFile bool) bool {
	if !env {
		if !arg {
			return fromFile
		} else {
			return arg
		}
	} else {
		return env
	}
}

func parseConfigFile(path string) (*ConfigFromFile, error) {
	if path == "" {
		return &ConfigFromFile{}, nil
	}

	f, err := os.ReadFile(path)

	if err != nil {
		if os.IsNotExist(err) {
			return &ConfigFromFile{}, errors.New("config file not found")
		}
		return &ConfigFromFile{}, err
	}

	configFromFile := ConfigFromFile{}

	err = json.Unmarshal(f, &configFromFile)
	return &configFromFile, err
}
