package service

import (
	"context"
	"fmt"
	"math/rand"
	"net/url"

	"github.com/kirillmashkov/shortener.git/internal/config"
	"github.com/kirillmashkov/shortener.git/internal/model"
)

type storeURL interface {
	AddURL(ctx context.Context, url string, keyURL string) error
	GetURL(ctx context.Context, keyURL string) (string, bool)
	AddBatchURL(ctx context.Context, shortOriginalURL []model.ShortOriginalURL) error
}

type Service struct {
	storage storeURL
	cfg     config.ServerConfig
}

func New(storage storeURL, config config.ServerConfig) *Service {
	return &Service{storage: storage, cfg: config}
}

func (s *Service) GetShortURL(ctx context.Context, originalURL *url.URL) (string, bool) {
	key := originalURL.Path[len("/"):]
	url, exist := s.storage.GetURL(ctx, key)

	if !exist {
		return "", false
	}

	return url, true
}

func (s *Service) ProcessURL(ctx context.Context, originalURL string) (string, bool) {
	keyURL := s.keyURL()
	shortURL := s.shortURL(keyURL)
	if err := s.storage.AddURL(ctx, originalURL, keyURL); err != nil {
		return "", false
	}
	return shortURL, true

}

func (s *Service) ProcessURLBatch(ctx context.Context, originalURLs []model.URLToShortBatchRequest) ([]model.ShortToURLBatchResponse, error) {
	var soURLs []model.ShortOriginalURL
	var results []model.ShortToURLBatchResponse

	for _, originalURL := range originalURLs {
		keyURL := s.keyURL()
		shortURL := s.shortURL(keyURL)
		soURLs = append(soURLs, model.ShortOriginalURL{Key: keyURL, OriginalURL: originalURL.OriginalURL})
		results = append(results, model.ShortToURLBatchResponse{CorrelationID: originalURL.CorrelationID, ShortURL: shortURL})
	}

	err := s.storage.AddBatchURL(ctx, soURLs)
	if err != nil {
		return nil, err
	}

	return results, nil
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
