package service

import (
	"fmt"
	"math/rand"
	"net/url"

	"github.com/kirillmashkov/shortener.git/internal/config"
)

type storeURL interface {
	AddURL(url string, keyURL string) error
	GetURL(keyURL string) (string, bool)
}

type Service struct {
	storage storeURL
	cfg     config.ServerConfig
}

func New(storage storeURL) *Service {
	return &Service{storage: storage}
}

func (s *Service) GetShortURL(originalURL *url.URL) (string, bool) {
	key := originalURL.Path[len("/"):]
	url, exist := s.storage.GetURL(key)

	if !exist {
		return "", false
	}

	return url, true
}

func (s *Service) ProcessURL(originalURL string) (string, bool) {
	keyURL := s.keyURL()
	shortURL := s.shortURL(keyURL)
	if err := s.storage.AddURL(originalURL, keyURL); err != nil {
		return "", false
	}
	return shortURL, true

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
