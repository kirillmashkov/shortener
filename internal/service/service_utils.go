package service

import (
	"context"

	"github.com/kirillmashkov/shortener.git/internal/storage/database"
	"go.uber.org/zap"
)

// ServiceUtils - тип для утильных функций
type ServiceUtils struct {
	db  *database.Database
	log *zap.Logger
}

// NewServiceUtils - конструктор
func NewServiceUtils(db *database.Database, log *zap.Logger) *ServiceUtils {
	return &ServiceUtils{db: db, log: log}
}

// PingDB - проверка доступнгсти БД
func (su *ServiceUtils) PingDB(ctx context.Context) error {
	if err := su.db.Ping(ctx); err != nil {
		su.log.Error("Error ping DB", zap.Error(err))
		return err
	}

	return nil
}
