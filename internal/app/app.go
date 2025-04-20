package app

import (
	// "github.com/kirillmashkov/shortener.git/internal/storage"
	"go.uber.org/zap"
)

type ServerConfig struct {
	Host        string "env:\"SERVER_ADDRESS\""
	Redirect    string "env:\"BASE_URL\""
	FileStorage string "env:\"FILE_STORAGE_PATH\""
}

var ServerEnv ServerConfig
var ServerArg ServerConfig
var ServerConf ServerConfig

// var StoreURL storage.StoreURLMap = *storage.InitStorage()

var Log *zap.Logger = zap.NewNop()
