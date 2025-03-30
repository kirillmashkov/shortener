package main

import (
	"fmt"
	"io"
	"math/rand"
	"net/http"
	"regexp"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

var urls = make(map[string]string)

func main() {
	r := chi.NewRouter()
	r.Get("/", postHandler)
	r.Post("/{id}", getHandler)
	r.Use(middleware.Logger)

	err := http.ListenAndServe(`:8080`, r)
	if err != nil {
		panic(err)
	}
}

func getHandler(res http.ResponseWriter, req *http.Request) {
	if req.Method != http.MethodGet {
		http.Error(res, "Only GET requests are allowed!", http.StatusBadRequest)
		return
	}

	key := req.URL.Path[len("/"):]
	url, exist := urls[key]

	if !exist {
		http.Error(res, "Key not found", http.StatusBadRequest)
		return
	}

	http.Redirect(res, req, url, http.StatusTemporaryRedirect)

}

func postHandler(res http.ResponseWriter, req *http.Request) {
	if req.Method != http.MethodPost {
		http.Error(res, "Only POST requests are allowed!", http.StatusBadRequest)
		return
	}

	originalURL, err := io.ReadAll(req.Body)
	if err != nil {
		http.Error(res, "Something went wrong", http.StatusBadRequest)
		return
	}

	validLink := validateLink(string(originalURL))
	
	if !validLink {
		http.Error(res, "Link is bad", http.StatusBadRequest)
		return
	}
		
	keyURL := keyURL()
	shortURL := shortURL(keyURL)
	urls[keyURL] = string(originalURL)
	res.Header().Set("content-type", "text/plain")
	res.WriteHeader(http.StatusCreated)
	res.Write([]byte(shortURL))
}

func validateLink(url string) bool {
	matched, _ := regexp.MatchString("^(http:\\/\\/www\\.|https:\\/\\/www\\.|http:\\/\\/|https:\\/\\/|\\/|\\/\\/)?[A-z0-9_-]*?[:]?[A-z0-9_-]*?[@]?[A-z0-9]+([\\-\\.]{1}[a-z0-9]+)*\\.[a-z]{2,5}(:[0-9]{1,5})?(\\/.*)?$", url)
	return matched
}

func keyURL() string {
	const dictionary = "ABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	const keyLen = 8

	keyURL := make([]byte, keyLen)
	for i := range keyURL {
		keyURL[i] = dictionary[rand.Intn(len(dictionary))]
	}
	return string(keyURL)
}

func shortURL(key string) string {
	return fmt.Sprintf("http://localhost:8080/%s", key)
}
