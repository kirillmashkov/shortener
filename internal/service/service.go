package service

import (
	"fmt"
	"math/rand"
	"net/url"
	"regexp"
	"strings"

	"github.com/kirillmashkov/shortener.git/internal/app"
	"github.com/kirillmashkov/shortener.git/internal/storage"
	"go.uber.org/zap"
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
	app.Log.Info("Original URL", zap.String("Original URL", originalURL))
	url := strings.ReplaceAll(string(originalURL), "\n", "")
	url = strings.ReplaceAll(string(url), "\\n", "")
	app.Log.Info("After trailing", zap.String("url", url))

	validLink := validateLink(url)
	if !validLink {
		return "", false
	}

	keyURL := keyURL()
	shortURL := shortURL(keyURL)
	storage.StoreURL.AddURL(url, keyURL)
	return shortURL, true

}

func validateLink(url string) bool {
	matched, _ := regexp.MatchString("^(http:\\/\\/www\\.|https:\\/\\/www\\.|http:\\/\\/|https:\\/\\/|\\/|\\/\\/)?[A-z0-9_-]*?[:]?[A-z0-9_-]*?[@]?[A-z0-9]+([\\-\\.]{1}[a-z0-9]+)*\\.[a-z]{2,5}(:[0-9]{1,5})?(\\/.*)?$", url)
	return matched
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
