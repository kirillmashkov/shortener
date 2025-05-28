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
	"github.com/kirillmashkov/shortener.git/internal/httpserver/middleware/security"
	"github.com/kirillmashkov/shortener.git/internal/model"

	"go.uber.org/zap"
)

type ServiceShortURL interface {
	GetShortURL(ctx context.Context, originalURL *url.URL) (string, bool, bool)
	ProcessURL(ctx context.Context, originalURL string, userID int) (string, error)
	ProcessURLBatch(ctx context.Context, originalURLs []model.URLToShortBatchRequest, userID int) ([]model.ShortToURLBatchResponse, error)
	DeleteURLBatch(userID int, shortURLs []string)
	GetAllURL(ctx context.Context, userID int) ([]model.ShortOriginalURL, error)
}

func GetHandler(res http.ResponseWriter, req *http.Request) {
	if req.Method != http.MethodGet {
		http.Error(res, "Only GET requests are allowed!", http.StatusBadRequest)
		return
	}

	url, deleted, exist := app.Service.GetShortURL(req.Context(), req.URL)

	if !exist {
		http.Error(res, "Key not found", http.StatusBadRequest)
		return
	}

	if deleted {
		res.WriteHeader(http.StatusGone)
		return
	}

	http.Redirect(res, req, url, http.StatusTemporaryRedirect)
}

func GetAllURL(res http.ResponseWriter, req *http.Request) {
	if req.Method != http.MethodGet {
		http.Error(res, "Only GET requests are allowed!", http.StatusBadRequest)
		return
	}

	u := security.UserIDType("userID")

	result, err := app.Service.GetAllURL(req.Context(), req.Context().Value(u).(int))
	if err != nil {
		http.Error(res, "Something went wrong", http.StatusBadRequest)
		return
	}

	res.Header().Set("Content-Type", "application/json")
	if len(result) == 0 {
		res.WriteHeader(http.StatusNoContent)
		return
	}

	res.WriteHeader(http.StatusOK)
	encoder := json.NewEncoder(res)
	if err := encoder.Encode(result); err != nil {
		app.Log.Debug("error encoding result", zap.Error(err))
		return
	}
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

	u := security.UserIDType("userID")
	
	shortURL, err := app.Service.ProcessURL(req.Context(), string(originalURL), req.Context().Value(u).(int))
	res.Header().Set("content-type", "text/plain")
	if err != nil {
		errorString := fmt.Sprintf("Something went wrong when generate short url for %s", string(originalURL))
		if errors.Is(err,model.ErrDuplicateURL) {
			res.WriteHeader(http.StatusConflict)
			_, err = res.Write([]byte(shortURL))
			if err != nil {
				app.Log.Error("Can't write response", zap.Error(err))
				http.Error(res, errorString, http.StatusBadRequest)
				return
			}
			return
		}
		app.Log.Error("Error process URL", zap.Error(err))
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

	u := security.UserIDType("userID")
	shortURL, err := app.Service.ProcessURL(req.Context(), request.OriginalURL, req.Context().Value(u).(int))

	res.Header().Set("Content-Type", "application/json")
	response := model.ShortToURLReponse{
		ShortURL: shortURL,
	}

	if err != nil {
		if errors.Is(err, model.ErrDuplicateURL) {
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

	u := security.UserIDType("userID")
	response, err := app.Service.ProcessURLBatch(req.Context(), request, req.Context().Value(u).(int))

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

func DeleteURLBatch(res http.ResponseWriter, req *http.Request) {
	if req.Method != http.MethodDelete {
		http.Error(res, "Only Delete requests are allowed!", http.StatusBadRequest)
		return
	}

	var shortURLs []string
	decoder := json.NewDecoder(req.Body)
	if err := decoder.Decode(&shortURLs); err != nil {
		app.Log.Debug("cannot parse request JSON body", zap.Error(err))
		http.Error(res, "cannot parse request JSON body", http.StatusBadRequest)
		return
	}

	u := security.UserIDType("userID")
	app.Service.DeleteURLBatch(req.Context().Value(u).(int), shortURLs)
	res.WriteHeader(http.StatusAccepted)
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
