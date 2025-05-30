package database

import (
	"context"
	"errors"
	"time"

	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/jackc/pgx/v5"
	"github.com/kirillmashkov/shortener.git/internal/config"
	"go.uber.org/zap"
	"github.com/jackc/pgx/v5/pgxpool"
)

const migrateDir = "migrations"
const timeoutPindDB = 1 * time.Second

type Database struct {
	cfg     *config.ServerConfig
	conn	*pgx.Conn
	dbpool	*pgxpool.Pool
	logger	*zap.Logger
}

func New(config *config.ServerConfig, logger *zap.Logger) *Database {
	return &Database{cfg: config, logger: logger}
}

func (d *Database) Ping(ctx context.Context) error {
	if d.conn == nil {
		return errors.New("no connection to db")
	}
	ctx, cancel := context.WithTimeout(ctx, timeoutPindDB)
	defer cancel()

	return d.conn.Ping(ctx)
}

func (d *Database) Open() error {
	var err error
	d.conn, err = pgx.Connect(context.Background(), d.cfg.Connection)
	d.dbpool, _ = pgxpool.New(context.Background(), d.cfg.Connection)
	return err
}

func (d *Database) Close() error {
	err :=  d.conn.Close(context.Background())
	d.dbpool.Close()
	return err
}

func (d *Database) Migrate() error {
	m, err := migrate.New("file://" + migrateDir, d.cfg.Connection)
	if err != nil {
		d.logger.Error("Can't initialize migrations", zap.Error(err))
		return err
	}

	if err := m.Up(); err != nil {
		if errors.Is(err, migrate.ErrNoChange) {
			d.logger.Info("No migrations need")
			return nil
		}
		d.logger.Error("Something went wrong while migrations", zap.Error(err))
		return err
	}
	return nil
}