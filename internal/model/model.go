package model

type URLToShortRequest struct {
	OriginalURL string `json:"url"` 
}

type ShortToURLReponse struct {
	ShortURL string `json:"result"`
}

type URLToShortBatchRequest struct {
	CorrelationId string `json:"correlation_id"`
	OriginalURL string `json:"original_url"`
}

type ShortToURLBatchResponse struct {
	CorrelationId string `json:"correlation_id"`
	ShortURL string `json:"short_url"`
}

type ShortOriginalUrl struct {
	Key string
	OriginalURL string
}
