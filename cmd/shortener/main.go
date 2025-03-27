package main

import (
	"fmt"
	"io"
	"math/rand"
	"net/http"
)

var urls = make(map[string]string)

func main() {
	mux := http.NewServeMux()
	mux.HandleFunc("/", postHandler)
	mux.HandleFunc("/{id}", getHandler)

	err := http.ListenAndServe(`:8080`, mux)
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
		http.Error(res, "Key not found", http.StatusBadRequest)
		return
	}

	keyURL := keyURL()
	shortURL := shortURL(keyURL)
	urls[keyURL] = string(originalURL)
	res.WriteHeader(http.StatusCreated)
	res.Write([]byte(shortURL))
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
