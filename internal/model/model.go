package model

import (
	"errors"
	"sync"
)

// URLToShortRequest - запрос с исходной ссылкой
type URLToShortRequest struct {
	OriginalURL string `json:"url"`
}

// URLToShortRequest - ответ с короткой ссылкой
type ShortToURLReponse struct {
	ShortURL string `json:"result"`
}

// URLToShortBatchRequest - запрос с исходной ссылкой и ключом корреляции
type URLToShortBatchRequest struct {
	CorrelationID string `json:"correlation_id"`
	OriginalURL   string `json:"original_url"`
}

// ShortToURLBatchResponse - ответ с короткой ссылкой и ключом корреляции
type ShortToURLBatchResponse struct {
	CorrelationID string `json:"correlation_id"`
	ShortURL      string `json:"short_url"`
}

// KeyOriginalURL - короткая ссылка + исходная ссылка
type KeyOriginalURL struct {
	Key         string
	OriginalURL string
}

// ShortOriginalURL - короткая ссылка + исходная ссылка
type ShortOriginalURL struct {
	Short       string `json:"short_url"`
	OriginalURL string `json:"original_url"`
}

// ErrDuplicateURL - ошибка дублирования url
var ErrDuplicateURL = errors.New("duplicate url")

// ShortURLUserID - для запроса на удаления ссылок для конкретного пользователя
type ShortURLUserID struct {
	ShortURLs []string
	UserID    int
}

// ShortURLchan - канал для асинхроннго удаления ссылок из БД
var ShortURLchan chan ShortURLUserID

// Wg - WaitGroup
var Wg *sync.WaitGroup
