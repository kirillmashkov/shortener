package app

import (
	"github.com/kirillmashkov/shortener.git/internal/storage"
)

type ServerConfig struct {
	Host     string "env:\"SERVER_ADDRESS\""
	Redirect string "env:\"BASE_URL\""
}

var ServerEnv ServerConfig
var ServerArg ServerConfig
var ServerConf ServerConfig

var StoreURL storage.StoreURLMap = *storage.NewStoreMap()