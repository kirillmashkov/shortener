package model

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

type ShortOriginalURL struct {
	Key         string `json:"short_url"`
	OriginalURL string `json:"original_url"`
}
