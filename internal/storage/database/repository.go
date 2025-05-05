package database

import (
	"context"
	"errors"
	"time"

	"github.com/google/uuid"
	"go.uber.org/zap"
)

type RepositoryShortURL struct {
	db *Database
	log *zap.Logger
}

func NewRepositoryShortURL(db *Database, log *zap.Logger) *RepositoryShortURL {
	return &RepositoryShortURL{db: db, log: log}
}

func (r *RepositoryShortURL) AddURL(ctx context.Context, url string, keyURL string) error {
	ctx, cancel := context.WithTimeout(ctx, 1 * time.Second)
	defer cancel()

	tx, err := r.db.conn.Begin(ctx)
	if err != nil {
		r.log.Error("Error open tran", zap.Error(err))
		return err
	}
	
	_, err = r.db.conn.Exec(ctx, "insert into shorturl (id, short_url, original_url) values ($1, $2, $3)", uuid.NewString(), keyURL, url)
	if err != nil {
		r.log.Error("Error insert short url ", 
			zap.String("key", keyURL),
			zap.String("original url", url), 
			zap.Error(err))
		
		if errRollback := tx.Rollback(ctx); errRollback != nil {
			return errors.Join(err, errRollback)
		}

		return err
	}

	err = tx.Commit(ctx)
	if err != nil {
		r.log.Error("Error commit tran", zap.Error(err))
	}
	
	return nil
}

func (r *RepositoryShortURL) GetURL(ctx context.Context, keyURL string) (string, bool) {
	ctx, cancel := context.WithTimeout(ctx, 1 * time.Second)
	defer cancel()

	var original_url string
	err := r.db.conn.QueryRow(ctx, "select original_url from shorturl where short_url = $1", keyURL).Scan(&original_url)
	if err != nil {
		return "", false
	}

	return original_url, true
}