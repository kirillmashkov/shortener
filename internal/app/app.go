package app

import (
	"github.com/kirillmashkov/shortener.git/internal/config"
	"github.com/kirillmashkov/shortener.git/internal/database"
	"github.com/kirillmashkov/shortener.git/internal/service"
	"github.com/kirillmashkov/shortener.git/internal/storage"
	"go.uber.org/zap"
)

var ServerConf config.ServerConfig

var Storage *storage.StoreURLMap
var Service *service.Service
var ServiceUtils *service.ServiceUtils

var Database *database.Database

var Log *zap.Logger = zap.NewNop()

func Initialize() error {
	var err error
	
	config.InitServerConf(&ServerConf, Log)

	Database = database.New(&ServerConf)
	err = Database.Open()
	if err != nil {
		Log.Error("Error open connection db", zap.Error(err))
	}

	Storage, err = storage.New(&ServerConf, Log, &ServerConf)
	if err != nil {
		return err
	}

	Service = service.New(Storage, ServerConf)
	ServiceUtils = service.NewServiceUtils(Database, Log)

	return nil
}

func Close() {
	
	errClose := Database.Close()
	if errClose != nil {
		Log.Error("Error close connection db", zap.Error(errClose))
	}
	
}