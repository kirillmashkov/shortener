package service

import (
	"fmt"
	"math/rand"
	"net/url"
	"github.com/kirillmashkov/shortener.git/internal/app"
	"github.com/kirillmashkov/shortener.git/internal/storage"
)

type storeURL interface {
	AddURL(url string, keyURL string)
	GetURL() (string, bool)
}

func GetShortURL(originalURL *url.URL) (string, bool) {
	key := originalURL.Path[len("/"):]
	url, exist := storage.StoreURL.GetURL(key)

	if !exist {
		return "", false
	}

	return url, true
}

func ProcessURL(originalURL string) (string, bool) {
	keyURL := keyURL()
	shortURL := shortURL(keyURL)
	storage.StoreURL.AddURL(originalURL, keyURL)
	return shortURL, true

}

func keyURL() string {
	const dictionary = "ABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	const keyLen = 8

	keyURL := make([]byte, keyLen)
	for i := range keyURL {
		keyURL[i] = dictionary[rand.Intn(len(dictionary))]
	}
	return string(keyURL)
}

func shortURL(key string) string {
	return fmt.Sprintf("%s/%s", app.ServerConf.Redirect, key)
}
