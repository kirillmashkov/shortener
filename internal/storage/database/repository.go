package database

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgerrcode"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/kirillmashkov/shortener.git/internal/model"
	"go.uber.org/zap"
)

type RepositoryShortURL struct {
	db  *Database
	log *zap.Logger
}

const timeoutOperationDB = 1 * time.Second

func NewRepositoryShortURL(db *Database, log *zap.Logger) *RepositoryShortURL {
	return &RepositoryShortURL{db: db, log: log}
}

func (r *RepositoryShortURL) AddURL(ctx context.Context, url string, keyURL string, userID int) error {
	ctx, cancel := context.WithTimeout(ctx, timeoutOperationDB)
	defer cancel()

	// tx, err := r.db.conn.Begin(ctx)
	tx, err := r.db.dbpool.Begin(ctx)
	if err != nil {
		r.log.Error("Error open tran", zap.Error(err))
		// return err
		return fmt.Errorf("error open tran %w", err)
	}
	defer func() {
		if err == nil {
			if errCommit := tx.Commit(ctx); errCommit != nil {
				r.log.Error("Error commit tran", zap.Error(err))
			}
		} else {
			if errRollback := tx.Rollback(ctx); errRollback != nil {
				r.log.Error("Error rollback tx", zap.Error(errRollback))
			}
		}
	}()

	err = r.insertShortURL(ctx, tx, keyURL, url, userID)
	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) {
			if pgErr.Code == pgerrcode.UniqueViolation {
				return model.NewDuplicateURLError(err)
			}
		}

		// return err
		return fmt.Errorf("error insert url %w", err)
	}

	return nil
}

func (r *RepositoryShortURL) AddBatchURL(ctx context.Context, shortOriginalURL []model.KeyOriginalURL, userID int) error {
	ctx, cancel := context.WithTimeout(ctx, timeoutOperationDB)
	defer cancel()

	tx, err := r.db.conn.Begin(ctx)
	if err != nil {
		r.log.Error("Error open tran", zap.Error(err))
		return err
	}
	defer func() {
		if err == nil {
			if errCommit := tx.Commit(ctx); errCommit != nil {
				r.log.Error("Error commit tran", zap.Error(err))
			}
		} else {
			if errRollback := tx.Rollback(ctx); errRollback != nil {
				r.log.Error("Error rollback tx", zap.Error(errRollback))
			}
		}
	}()

	for _, soURL := range shortOriginalURL {
		err = r.insertShortURL(ctx, tx, soURL.Key, soURL.OriginalURL, userID)
		if err != nil {
			return err
		}
	}

	return nil
}

func (r *RepositoryShortURL) DeleteURLBatch(ctx context.Context, shortURL []string, userID int) error {
	ctx, cancel := context.WithTimeout(ctx, timeoutOperationDB)
	defer cancel()

	// tx, err := r.db.conn.Begin(ctx)
	tx, err := r.db.dbpool.Begin(ctx)
	if err != nil {
		r.log.Error("Error open tran", zap.Error(err))
		return err
	}
	defer func() {
		if err == nil {
			if errCommit := tx.Commit(ctx); errCommit != nil {
				r.log.Error("Error commit tran", zap.Error(err))
			}
		} else {
			if errRollback := tx.Rollback(ctx); errRollback != nil {
				r.log.Error("Error rollback tx", zap.Error(errRollback))
			}
		}
	}()

	batch := &pgx.Batch{}
	for _, data := range shortURL {
		r.log.Info("Set deleted = true", zap.String("short_url", data), zap.Int("userID", userID))
		batch.Queue("update shorturl set deleted = true where short_url = $1 and user_id = $2", data, userID)
	}

	res := tx.SendBatch(ctx, batch)
	
	errClose := res.Close()
	if errClose != nil {
		r.log.Error("Error close batch delete short url", zap.Error(err))
	}

	return nil
}

func (r *RepositoryShortURL) insertShortURL(ctx context.Context, tx pgx.Tx, keyURL string, url string, userID int) error {
	_, err := tx.Exec(ctx, "insert into shorturl (id, short_url, original_url, user_id) values ($1, $2, $3, $4)", uuid.NewString(), keyURL, url, userID)
	if err != nil {
		r.log.Error("Error insert short url ",
			zap.String("key", keyURL),
			zap.String("original url", url),
			zap.Error(err))
		return err
	}

	return nil
}

func (r *RepositoryShortURL) GetURL(ctx context.Context, keyURL string) (string, bool, bool) {
	ctx, cancel := context.WithTimeout(ctx, 1*time.Second)
	defer cancel()

	var originalURL string
	var deleted bool
	// err := r.db.conn.QueryRow(ctx, "select original_url, deleted from shorturl where short_url = $1", keyURL).Scan(&originalURL, &deleted)
	err := r.db.dbpool.QueryRow(ctx, "select original_url, deleted from shorturl where short_url = $1", keyURL).Scan(&originalURL, &deleted)
	if err != nil {
		r.log.Error("Error get originalUrl from db", zap.String("shortUrl", keyURL), zap.Error(err))
		return "", false, false
	}

	return originalURL, deleted, true
}

func (r *RepositoryShortURL) GetShortURL(ctx context.Context, originalURL string) (string, error) {
	ctx, cancel := context.WithTimeout(ctx, timeoutOperationDB)
	defer cancel()

	var key string
	err := r.db.conn.QueryRow(ctx, "select short_url from shorturl where original_url = $1", originalURL).Scan(&key)
	if err != nil {
		return "", err
	}

	return key, nil
}

func (r *RepositoryShortURL) GetAllURL(ctx context.Context, userID int) ([]model.KeyOriginalURL, error) {
	ctx, cancel := context.WithTimeout(ctx, timeoutOperationDB)
	defer cancel()

	rows, err := r.db.conn.Query(ctx, "select short_url, original_url from shorturl where user_id = $1", userID)
	if err != nil {
		r.log.Error("Error get all urls from db", zap.Error(err))
		return nil, err
	}

	defer rows.Close()

	res, err := pgx.CollectRows(rows, pgx.RowToStructByPos[model.KeyOriginalURL])
	for _, j := range res {
		r.log.Info("Row", zap.String("key", j.Key), zap.String("original", j.OriginalURL))
	}

	return res, err
}
