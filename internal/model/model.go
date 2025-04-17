package model

type URLToShortRequest struct {
	OriginalURL string `json:"url"` 
}

type ShortToURLReponse struct {
	ShortURL string `json:"result"`
}
