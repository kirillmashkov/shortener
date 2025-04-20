package storage

import (
	"bufio"
	"encoding/json"
	"os"
	"strconv"
	"sync"

	"github.com/kirillmashkov/shortener.git/internal/app"
	"go.uber.org/zap"
)

type StoreURLMap struct {
	sync.RWMutex
	urls map[string]string
}

type StoreFile struct {
	UUID        string `json:"uuid"`
	ShortURL    string `json:"short_url"`
	OriginalURL string `json:"original_url"`
}

var StoreURL StoreURLMap
var id int = 0

func InitStorage() {
	StoreURL.urls = map[string]string{}

	app.Log.Info("Read storage file", zap.String("file", app.ServerConf.FileStorage))
	file, err := os.OpenFile(app.ServerConf.FileStorage, os.O_RDONLY|os.O_CREATE, 0666)
	if err != nil {
		return
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		shortURL := StoreFile{}
		json.Unmarshal(scanner.Bytes(), &shortURL)
		id, err = strconv.Atoi(shortURL.UUID)
		if err == nil {
			app.Log.Info("Read short ulr", 
				zap.String("shortURL", shortURL.ShortURL),
				zap.String("OriginalURL", shortURL.OriginalURL),
				zap.Int("id", id))
			StoreURL.urls[shortURL.ShortURL] = shortURL.OriginalURL
		}
	}

	if err := scanner.Err(); err != nil {
		app.Log.Error("Error read file", zap.Error(err))
	}
	id++
}

func (storeMap *StoreURLMap) AddURL(url string, keyURL string) {
	storeMap.Lock()
	err := saveShortURLToFile(keyURL, url)
	if err == nil {
		storeMap.urls[keyURL] = url
		id++
	}
	storeMap.Unlock()
}

func (storeMap *StoreURLMap) GetURL(keyURL string) (string, bool) {
	storeMap.RLock()
	url, exist := storeMap.urls[keyURL]
	storeMap.RUnlock()
	return url, exist
}

func saveShortURLToFile(url string, originalURL string) error {
	app.Log.Info("Write to file storage", zap.String("file", app.ServerConf.FileStorage))
	file, err := os.OpenFile(app.ServerConf.FileStorage, os.O_WRONLY|os.O_APPEND|os.O_CREATE, 0666)
	if err != nil {
		return err
	}

	shortURLToFile := StoreFile{
		UUID:        strconv.Itoa(id),
		ShortURL:    url,
		OriginalURL: originalURL,
	}

	writer := bufio.NewWriter(file)

	app.Log.Info("Write short url", zap.Any("short url", shortURLToFile))
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

	writer.Flush()
	file.Close()

	return nil
}
