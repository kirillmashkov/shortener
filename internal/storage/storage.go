package storage

import (
	"sync"
)

type StoreURL interface {
	AddURL(url string, keyURL string)
	GetURL() (string, bool)
}

type StoreURLMap struct {
	sync.RWMutex
	urls map[string]string
}

func NewStoreMap() *StoreURLMap {
	var storeMap StoreURLMap
	storeMap.urls = map[string]string{}
	return &storeMap
}

func (storeMap *StoreURLMap) AddURL(url string, keyURL string) {
	storeMap.Lock()
	storeMap.urls[keyURL] = url
	storeMap.Unlock()
}

func (storeMap *StoreURLMap) GetURL(keyURL string) (string, bool) {
	storeMap.RLock()
	url, exist := storeMap.urls[keyURL]
	storeMap.RUnlock()
	return url, exist
}
