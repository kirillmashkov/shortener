package handler

import (
	"net/http"
	"github.com/kirillmashkov/shortener.git/internal/service"
	"io"
	"fmt"
)

func GetHandler(res http.ResponseWriter, req *http.Request) {
	if req.Method != http.MethodGet {
		http.Error(res, "Only GET requests are allowed!", http.StatusBadRequest)
		return
	}
	
	url, exist := service.GetShortURL(req.URL)

	if !exist {
		http.Error(res, "Key not found", http.StatusBadRequest)
		return
	}

	http.Redirect(res, req, url, http.StatusTemporaryRedirect)
}

func PostHandler(res http.ResponseWriter, req *http.Request) {
	if req.Method != http.MethodPost {
		http.Error(res, "Only POST requests are allowed!", http.StatusBadRequest)
		return
	}

	originalURL, err := io.ReadAll(req.Body)
	if err != nil {
		http.Error(res, "Something went wrong", http.StatusBadRequest)
		return
	}

	shortURL, result := service.ProcessURL(string(originalURL))
	if !result {
		errorString := fmt.Sprintf("Link is bad %s", string(originalURL))
		http.Error(res, errorString, http.StatusBadRequest)
		return
	}

	res.Header().Set("content-type", "text/plain")
	res.WriteHeader(http.StatusCreated)
	res.Write([]byte(shortURL))

}