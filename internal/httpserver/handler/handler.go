package handler

import (
	"context"
	"encoding/json"
	"errors"
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
	ProcessURL(ctx context.Context, originalURL string) (string, error)
	ProcessURLBatch(ctx context.Context, originalURLs []model.URLToShortBatchRequest) ([]model.ShortToURLBatchResponse, error)
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

	shortURL, err := app.Service.ProcessURL(req.Context(), string(originalURL))
	res.Header().Set("content-type", "text/plain")
	if err != nil {
		var errDuplicate *model.DuplicateURLError
		if errors.As(err, &errDuplicate) {
			res.WriteHeader(http.StatusConflict)	
			_, err = res.Write([]byte(shortURL))
			if err != nil {
				http.Error(res, "Can't write response", http.StatusBadRequest)
				return
			}
			return
		}
		errorString := fmt.Sprintf("Something went wrong when generate short url for %s", string(originalURL))
		http.Error(res, errorString, http.StatusBadRequest)
		return
	}

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

	shortURL, err := app.Service.ProcessURL(req.Context(), request.OriginalURL)

	res.Header().Set("Content-Type", "application/json")
	response := model.ShortToURLReponse {
		ShortURL: shortURL,
	}

	if err != nil {
		var errDuplicate *model.DuplicateURLError
		if errors.As(err, &errDuplicate) {
			res.WriteHeader(http.StatusConflict)	
			encoder := json.NewEncoder(res)
			if err := encoder.Encode(response); err != nil {
				app.Log.Debug("error encoding response", zap.Error(err))
				return
			}
			return
		}

		errorString := fmt.Sprintf("Something went wrong when generate short url for %s", request.OriginalURL)
		http.Error(res, errorString, http.StatusBadRequest)
		return
	}

	res.WriteHeader(http.StatusCreated)
	encoder := json.NewEncoder(res)
    if err := encoder.Encode(response); err != nil {
        app.Log.Debug("error encoding response", zap.Error(err))
        return
    }
}

func PostGenerateShortURLBatch(res http.ResponseWriter, req *http.Request) {
	if req.Method != http.MethodPost {
		http.Error(res, "Only POST requests are allowed!", http.StatusBadRequest)
		return
	}

	var request []model.URLToShortBatchRequest
	decoder := json.NewDecoder(req.Body)
	if err := decoder.Decode(&request); err != nil {
		app.Log.Debug("cannot parse request JSON body", zap.Error(err))
		http.Error(res, "cannot parse request JSON body", http.StatusBadRequest)
		return
	}

	response, err := app.Service.ProcessURLBatch(req.Context(), request)
	if err != nil {
		http.Error(res, "Can't store url batch", http.StatusBadRequest)
		return
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

	err := app.ServiceUtils.PingDB(req.Context())
	if err != nil {
		http.Error(res, "DB is unavailable", http.StatusInternalServerError)
		return
	}

	res.WriteHeader(http.StatusOK)
}
