// Модуль service - слой бизнес логики по управлению короткими ссылками
package service

import (
	"context"
	"errors"
	"fmt"
	"math/rand"
	"net/url"

	"github.com/kirillmashkov/shortener.git/internal/config"
	"github.com/kirillmashkov/shortener.git/internal/model"
	"go.uber.org/zap"
)

type storeURL interface {
	AddURL(ctx context.Context, url string, keyURL string, userID int) error
	GetURL(ctx context.Context, keyURL string) (string, bool, bool)
	GetAllURL(ctx context.Context, userID int) ([]model.KeyOriginalURL, error)
	AddBatchURL(ctx context.Context, shortOriginalURL []model.KeyOriginalURL, userID int) error
	DeleteURLBatchProcessor(ctx context.Context)
	GetShortURL(ctx context.Context, originalURL string) (string, error)
	GetStats(ctx context.Context) (int, int, error)
}

// Service - тип для сервисного слоя по управлению ссылками
type Service struct {
	storage storeURL
	cfg     config.ServerConfig
	log     *zap.Logger
}

// New - конструктор
func New(storage storeURL, config config.ServerConfig, log *zap.Logger) *Service {
	return &Service{storage: storage, cfg: config, log: log}
}

// GetShortURL возвращает исходную ссылку по короткому названию
func (s *Service) GetShortURL(ctx context.Context, originalURL *url.URL) (string, bool, bool) {
	key := originalURL.Path[len("/"):]
	url, deleted, exist := s.storage.GetURL(ctx, key)

	if !exist {
		return "", false, false
	}

	if deleted {
		return "", true, true
	}

	return url, false, true
}

// GetAllURL - возвращает все ссылки для пользователя
func (s *Service) GetAllURL(ctx context.Context, userID int) ([]model.ShortOriginalURL, error) {
	keyShortURL, err := s.storage.GetAllURL(ctx, userID)
	if err != nil {
		return nil, err
	}

	result := make([]model.ShortOriginalURL, 0, len(keyShortURL))
	for _, j := range keyShortURL {
		result = append(result, model.ShortOriginalURL{Short: s.shortURL(j.Key), OriginalURL: j.OriginalURL})
	}

	return result, nil
}

// ProcessURL - сохраняет исходную ссылку, возвращает короткую ссылку
func (s *Service) ProcessURL(ctx context.Context, originalURL string, userID int) (string, error) {
	keyURL := s.keyURL()
	shortURL := s.shortURL(keyURL)
	if err := s.storage.AddURL(ctx, originalURL, keyURL, userID); err != nil {
		if errors.Is(err, model.ErrDuplicateURL) {
			key, errGetShortURL := s.storage.GetShortURL(ctx, originalURL)
			if errGetShortURL != nil {
				return "", errors.New("can't get short url")
			}

			return s.shortURL(key), err
		}
		return "", err
	}
	return shortURL, nil
}

// ProcessURLBatch - сохранение массива ссылок, возвращает ключ и короткую ссылку
func (s *Service) ProcessURLBatch(ctx context.Context, originalURLs []model.URLToShortBatchRequest, userID int) ([]model.ShortToURLBatchResponse, error) {
	var soURLs []model.KeyOriginalURL
	var results []model.ShortToURLBatchResponse

	for _, originalURL := range originalURLs {
		keyURL := s.keyURL()
		shortURL := s.shortURL(keyURL)
		soURLs = append(soURLs, model.KeyOriginalURL{Key: keyURL, OriginalURL: originalURL.OriginalURL})
		results = append(results, model.ShortToURLBatchResponse{CorrelationID: originalURL.CorrelationID, ShortURL: shortURL})
	}

	err := s.storage.AddBatchURL(ctx, soURLs, userID)
	if err != nil {
		return nil, err
	}

	return results, nil
}

// DeleteURLBatch - удаление массива ссылок для пользователя
func (s *Service) DeleteURLBatch(userID int, shortURLs []string) {
	shortURLUser := model.ShortURLUserID{ShortURLs: shortURLs, UserID: userID}
	model.ShortURLchan <- shortURLUser
}

func (s *Service) keyURL() string {
	const dictionary = "ABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	const keyLen = 8

	keyURL := make([]byte, keyLen)
	for i := range keyURL {
		keyURL[i] = dictionary[rand.Intn(len(dictionary))]
	}
	return string(keyURL)
}

func (s *Service) shortURL(key string) string {
	return fmt.Sprintf("%s/%s", s.cfg.Redirect, key)
}

func (s *Service) GetStats(ctx context.Context) (model.Stats, error) {
	usersCount, urlsCount, err := s.storage.GetStats(ctx)
	return model.Stats{UrlsCount: urlsCount, UsersCount: usersCount}, err
}
