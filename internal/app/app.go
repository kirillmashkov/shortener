package app

import (
	"github.com/kirillmashkov/shortener.git/internal/storage"
	"github.com/kirillmashkov/shortener.git/internal/service"
	"github.com/kirillmashkov/shortener.git/internal/config"
	"go.uber.org/zap"
)

var ServerConf config.ServerConfig

var Storage *storage.StoreURLMap
var Service *service.Service

var Log *zap.Logger = zap.NewNop()

func Initialize() error {
	var err error
	
	config.InitServerConf(&ServerConf, Log)

	Storage, err = storage.New(&ServerConf, Log, &ServerConf)
	if err != nil {
		return err
	}

	Service = service.New(Storage) 

	return nil
}
