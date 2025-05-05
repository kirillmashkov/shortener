package handler

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"

	"github.com/kirillmashkov/shortener.git/internal/app"
	"github.com/kirillmashkov/shortener.git/internal/model"
	"go.uber.org/zap"
)

type ServiceShortURL interface {
	GetShortURL(ctx context.Context, originalURL *url.URL) (string, bool)
	ProcessURL(ctx context.Context, originalURL string) (string, bool)
}

func GetHandler(res http.ResponseWriter, req *http.Request) {
	if req.Method != http.MethodGet {
		http.Error(res, "Only GET requests are allowed!", http.StatusBadRequest)
		return
	}

	url, exist := app.Service.GetShortURL(req.Context(), req.URL)

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


	shortURL, result := app.Service.ProcessURL(req.Context(), string(originalURL))
	if !result {
		errorString := fmt.Sprintf("Link is bad %s", string(originalURL))
		http.Error(res, errorString, http.StatusBadRequest)
		return
	}

	res.Header().Set("content-type", "text/plain")
	res.WriteHeader(http.StatusCreated)
	_, err = res.Write([]byte(shortURL))
	if err != nil {
		http.Error(res, "Can't write response", http.StatusBadRequest)
		return
	}
}

func PostGenerateShortURL(res http.ResponseWriter, req *http.Request) {
	if req.Method != http.MethodPost {
		http.Error(res, "Only POST requests are allowed!", http.StatusBadRequest)
		return
	}
	var request model.URLToShortRequest
	decoder := json.NewDecoder(req.Body)
	if err := decoder.Decode(&request); err != nil {
		
		app.Log.Debug("cannot parse request JSON body", zap.Error(err))
		http.Error(res, "cannot parse request JSON body", http.StatusBadRequest)
		return
	}

	shortURL, result := app.Service.ProcessURL(req.Context(), request.OriginalURL)
	if !result {
		errorString := fmt.Sprintf("Link is bad %s", string(request.OriginalURL))
		http.Error(res, errorString, http.StatusBadRequest)
		return
	}

	response := model.ShortToURLReponse {
		ShortURL: shortURL,
	}

	res.Header().Set("Content-Type", "application/json")
	res.WriteHeader(http.StatusCreated)
	encoder := json.NewEncoder(res)
    if err := encoder.Encode(response); err != nil {
        app.Log.Debug("error encoding response", zap.Error(err))
        return
    }
}

func Ping(res http.ResponseWriter, req *http.Request) {
	if req.Method != http.MethodGet {
		http.Error(res, "Only GET requests are allowed!", http.StatusBadRequest)
		return
	}

	err := app.ServiceUtils.PingDB()
	if err != nil {
		http.Error(res, "DB is unavailable", http.StatusInternalServerError)
		return
	}

	res.WriteHeader(http.StatusOK)
}
