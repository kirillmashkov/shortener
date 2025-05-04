package service

import (
	"github.com/kirillmashkov/shortener.git/internal/database"
	"go.uber.org/zap"
)

type ServiceUtils struct {
	db *database.Database
	log *zap.Logger
}

func NewServiceUtils(db *database.Database, log *zap.Logger) *ServiceUtils {
	return &ServiceUtils{db: db, log: log}
}

func (su *ServiceUtils) PingDB() error {
	if err := su.db.Ping(); err != nil {
		su.log.Error("Error ping DB", zap.Error(err))
		return err
	}

	return nil
}