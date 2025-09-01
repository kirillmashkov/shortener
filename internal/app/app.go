// Модуль инициализации, содержит общие переменные и их инициализацию
package app

import (
	"context"
	"sync"

	"github.com/kirillmashkov/shortener.git/internal/config"
	"github.com/kirillmashkov/shortener.git/internal/model"
	"github.com/kirillmashkov/shortener.git/internal/service"
	"github.com/kirillmashkov/shortener.git/internal/storage/database"
	"github.com/kirillmashkov/shortener.git/internal/storage/memory"
	"go.uber.org/zap"
)

// ServerConf - конфигурация приложения
var ServerConf config.ServerConfig

// Storage - управление хранением ссылок в памяти
var Storage *memory.StoreURLMap

// Service - управление ссылками
var Service *service.Service

// ServiceUtils - утильные функции
var ServiceUtils *service.ServiceUtils

// Database - управление подключением и миграцией в БД
var Database *database.Database

// repositoryShortURL - управление хранением ссылок в БД
var repositoryShortURL *database.RepositoryShortURL

// Log - логер
var Log *zap.Logger = zap.NewNop()

// Initialize - инициализация приложения
func Initialize(ctx context.Context) error {
	var err error

	config.InitServerConf(&ServerConf, Log)

	model.Wg = &sync.WaitGroup{}

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
		repositoryShortURL = database.NewRepositoryShortURL(Database, Log)
		Service = service.New(repositoryShortURL, ServerConf, Log)
		model.ShortURLchan = make(chan model.ShortURLUserID)
		model.Wg.Add(1)
		go repositoryShortURL.DeleteURLBatchProcessor(ctx)
	}
	ServiceUtils = service.NewServiceUtils(Database, Log)
	return nil
}

// Close - закрытие приложения
func Close() {
	errClose := Database.Close()
	if errClose != nil {
		Log.Error("Error close connection db", zap.Error(errClose))
	}
}
