package memory

import (
	"bufio"
	"context"
	"encoding/json"
	"errors"
	"os"
	"sync"

	"github.com/google/uuid"
	"github.com/kirillmashkov/shortener.git/internal/config"
	"github.com/kirillmashkov/shortener.git/internal/model"
	"go.uber.org/zap"
)

// StoreURLMap - доступ к хранения в памяти ссылок
type StoreURLMap struct {
	mu     sync.RWMutex
	urls   map[string]string
	logger *zap.Logger
	cfg    *config.ServerConfig
}

// StoreFile - json для сохранения ссылок в файл
type StoreFile struct {
	UUID        string `json:"uuid"`
	ShortURL    string `json:"short_url"`
	OriginalURL string `json:"original_url"`
}

// New - конструктор
func New(conf *config.ServerConfig, logger *zap.Logger, config *config.ServerConfig) (*StoreURLMap, error) {
	urls := map[string]string{}

	logger.Info("Read storage file", zap.String("file", conf.FileStorage))
	file, err := os.OpenFile(conf.FileStorage, os.O_RDONLY|os.O_CREATE, 0666)
	if err != nil {
		return nil, err
	}
	defer func() {
		if errClose := file.Close(); errClose != nil {
			logger.Error("Can't close storage file when read")
			err = errors.Join(errClose)
		}
	}()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		shortURL := StoreFile{}
		err = json.Unmarshal(scanner.Bytes(), &shortURL)
		if err != nil {
			logger.Error("Can't parse storage file")
			return nil, err
		}

		logger.Info("Read short ulr",
			zap.String("shortURL", shortURL.ShortURL),
			zap.String("OriginalURL", shortURL.OriginalURL),
			zap.String("id", shortURL.UUID))
		urls[shortURL.ShortURL] = shortURL.OriginalURL

	}

	if err := scanner.Err(); err != nil {
		logger.Error("Error read file", zap.Error(err))
	}

	return &StoreURLMap{
		urls:   urls,
		logger: logger,
		cfg:    config,
	}, nil
}

// AddURL - сохранение ссылки
func (storeMap *StoreURLMap) AddURL(ctx context.Context, url string, keyURL string, userID int) error {
	storeMap.mu.Lock()
	defer storeMap.mu.Unlock()

	err := storeMap.saveShortURLToFile(keyURL, url)
	if err != nil {
		storeMap.logger.Error("Can't save link into file")
		return err
	}

	storeMap.urls[keyURL] = url
	return nil
}

// AddBatchURL - сохранение массива ссылок
func (storeMap *StoreURLMap) AddBatchURL(ctx context.Context, shortOriginalURL []model.KeyOriginalURL, userID int) error {
	storeMap.mu.Lock()
	defer storeMap.mu.Unlock()
	err := storeMap.saveShortURLToFileBatch(shortOriginalURL)
	if err != nil {
		storeMap.logger.Error("Can't save links into file")
		return err
	}

	for _, soURL := range shortOriginalURL {
		storeMap.urls[soURL.Key] = soURL.OriginalURL
	}

	return nil
}

// GetURL - получение ссылки
func (storeMap *StoreURLMap) GetURL(ctx context.Context, keyURL string) (string, bool, bool) {
	storeMap.mu.RLock()
	url, exist := storeMap.urls[keyURL]
	storeMap.mu.RUnlock()
	return url, false, exist
}

// GetAllURL - получение всех ссылок
func (storeMap *StoreURLMap) GetAllURL(ctx context.Context, userID int) ([]model.KeyOriginalURL, error) {
	storeMap.mu.Lock()
	defer storeMap.mu.Unlock()

	res := make([]model.KeyOriginalURL, 0, len(storeMap.urls))
	for k, v := range storeMap.urls {
		res = append(res, model.KeyOriginalURL{Key: k, OriginalURL: v})
	}
	return res, nil
}

// GetShortURL - реализация отсутствует
func (storeMap *StoreURLMap) GetShortURL(ctx context.Context, originalURL string) (string, error) {
	return "", errors.New("unsupport operation")
}

// DeleteURLBatchProcessor - реализация отсутствует
func (storeMap *StoreURLMap) DeleteURLBatchProcessor() {
	storeMap.logger.Error("unsupport operation")
}

func (storeMap *StoreURLMap) saveShortURLToFileBatch(shortOriginalURL []model.KeyOriginalURL) error {
	storeMap.logger.Info("Write to file storage", zap.String("file", storeMap.cfg.FileStorage))
	file, err := os.OpenFile(storeMap.cfg.FileStorage, os.O_WRONLY|os.O_APPEND|os.O_CREATE, 0666)
	if err != nil {
		return err
	}
	defer func() {
		if errClose := file.Close(); errClose != nil {
			storeMap.logger.Error("Can't close storage file when save it")
		}
	}()

	writer := bufio.NewWriter(file)

	for _, soURL := range shortOriginalURL {
		err = storeMap.writeToFile(soURL.Key, soURL.OriginalURL, writer)
		if err != nil {
			return err
		}
	}

	if err := writer.Flush(); err != nil {
		return err
	}

	return nil
}

func (storeMap *StoreURLMap) saveShortURLToFile(url string, originalURL string) error {
	storeMap.logger.Info("Write to file storage", zap.String("file", storeMap.cfg.FileStorage))
	file, err := os.OpenFile(storeMap.cfg.FileStorage, os.O_WRONLY|os.O_APPEND|os.O_CREATE, 0666)
	if err != nil {
		return err
	}
	defer func() {
		if errClose := file.Close(); errClose != nil {
			storeMap.logger.Error("Can't close storage file when save it")
		}
	}()

	writer := bufio.NewWriter(file)

	if err := storeMap.writeToFile(url, originalURL, writer); err != nil {
		return err
	}

	if err := writer.Flush(); err != nil {
		return err
	}

	return nil
}

func (storeMap *StoreURLMap) writeToFile(shortURL string, originalURL string, writer *bufio.Writer) error {
	shortURLToFile := StoreFile{
		UUID:        uuid.NewString(),
		ShortURL:    shortURL,
		OriginalURL: originalURL,
	}

	storeMap.logger.Info("Write short url", zap.Any("short url", shortURLToFile))
	data, err := json.Marshal(shortURLToFile)
	if err != nil {
		return err
	}

	if _, err := writer.Write(data); err != nil {
		return err
	}

	if err := writer.WriteByte('\n'); err != nil {
		return err
	}

	return nil
}
