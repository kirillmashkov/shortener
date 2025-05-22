package app

import (
	"github.com/kirillmashkov/shortener.git/internal/config"
	"github.com/kirillmashkov/shortener.git/internal/storage/database"
	"github.com/kirillmashkov/shortener.git/internal/service"
	"github.com/kirillmashkov/shortener.git/internal/storage/memory"
	"go.uber.org/zap"
)

var ServerConf config.ServerConfig

var Storage *memory.StoreURLMap
var Service *service.Service
var ServiceUtils *service.ServiceUtils

var Database *database.Database
var RepositoryShortURL *database.RepositoryShortURL

var Log *zap.Logger = zap.NewNop()

func Initialize() error {
	var err error
	
	config.InitServerConf(&ServerConf, Log)	

	Database = database.New(&ServerConf, Log)
	err = Database.Open()
	if err != nil {
		Log.Error("Error open connection db", zap.Error(err))
		Storage, err = memory.New(&ServerConf, Log, &ServerConf)
		if err != nil {
			return nil
		}
		Service = service.New(Storage, ServerConf, Log)
	} else {
		if err := Database.Migrate(); err != nil {
			return err
		}
		RepositoryShortURL = database.NewRepositoryShortURL(Database, Log)
		Service = service.New(RepositoryShortURL, ServerConf, Log)
	}

	ServiceUtils = service.NewServiceUtils(Database, Log)

	return nil
}

func Close() {
	
	errClose := Database.Close()
	if errClose != nil {
		Log.Error("Error close connection db", zap.Error(errClose))
	}
	
}