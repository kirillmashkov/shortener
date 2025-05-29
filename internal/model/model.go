package model

import (
	"errors"
	"sync"
)

type URLToShortRequest struct {
	OriginalURL string `json:"url"`
}

type ShortToURLReponse struct {
	ShortURL string `json:"result"`
}

type URLToShortBatchRequest struct {
	CorrelationID string `json:"correlation_id"`
	OriginalURL   string `json:"original_url"`
}

type ShortToURLBatchResponse struct {
	CorrelationID string `json:"correlation_id"`
	ShortURL      string `json:"short_url"`
}

type KeyOriginalURL struct {
	Key         string
	OriginalURL string
}

type ShortOriginalURL struct {
	Short       string `json:"short_url"`
	OriginalURL string `json:"original_url"`
}

var ErrDuplicateURL = errors.New("duplicate url")

type ShortURLUserID struct {
	ShortURLs []string
	UserID int
}

var ShortURLchan chan ShortURLUserID

var Wg *sync.WaitGroup