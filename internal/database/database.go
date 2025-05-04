package database

import (
	"context"
	"errors"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/kirillmashkov/shortener.git/internal/config"
)

type Database struct {
	cfg     *config.ServerConfig
	conn	*pgx.Conn
}

func New(config *config.ServerConfig) *Database {
	return &Database{cfg: config}
}

func (d *Database) Ping() error {
	if d.conn == nil {
		return errors.New("no connection to db")
	}
	ctx, cancel := context.WithTimeout(context.Background(), 1 * time.Second)
	defer cancel()

	return d.conn.Ping(ctx)
}

func (d *Database) Open() error {
	var err error
	d.conn, err = pgx.Connect(context.Background(), d.cfg.Connection)
	return err
}

func (d *Database) Close() error {
	return d.conn.Close(context.Background())
}