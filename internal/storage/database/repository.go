package database

import (
	"context"
	"errors"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/kirillmashkov/shortener.git/internal/model"
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
	
	err = r.insertShortURL(ctx, tx, keyURL, url)
	if err != nil {
		return err
	}

	err = tx.Commit(ctx)
	if err != nil {
		r.log.Error("Error commit tran", zap.Error(err))
	}
	
	return nil
}

func (r *RepositoryShortURL) AddBatchURL(ctx context.Context, shortOriginalURL []model.ShortOriginalUrl) error {
	ctx, cancel := context.WithTimeout(ctx, 1 * time.Second)
	defer cancel()

	tx, err := r.db.conn.Begin(ctx)
	if err != nil {
		r.log.Error("Error open tran", zap.Error(err))
		return err
	}

	for _, soURL := range shortOriginalURL {
		err = r.insertShortURL(ctx, tx, soURL.Key, soURL.OriginalURL)
		if err != nil {
			return err
		}
	}

	err = tx.Commit(ctx)
	if err != nil {
		r.log.Error("Error commit tran", zap.Error(err))
	}

	return nil
}

func (r *RepositoryShortURL) insertShortURL(ctx context.Context, tx pgx.Tx, keyURL string, url string) error {
	_, err := tx.Exec(ctx, "insert into shorturl (id, short_url, original_url) values ($1, $2, $3)", uuid.NewString(), keyURL, url)
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

	return nil
}

func (r *RepositoryShortURL) GetURL(ctx context.Context, keyURL string) (string, bool) {
	ctx, cancel := context.WithTimeout(ctx, 1 * time.Second)
	defer cancel()

	var originalURL string
	err := r.db.conn.QueryRow(ctx, "select original_url from shorturl where short_url = $1", keyURL).Scan(&originalURL)
	if err != nil {
		return "", false
	}

	return originalURL, true
}