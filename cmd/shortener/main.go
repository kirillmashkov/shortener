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
		http.Error(res, "Only GET requests are allowed!", http.StatusMethodNotAllowed)
		return
	}

	key := req.URL.Path[len("/"):]
	url, exist := urls[key]

	if !exist {
		http.Error(res, "Key not found", http.StatusBadRequest)
		return
	}

	http.Redirect(res, req, url, http.StatusMovedPermanently)

}

func postHandler(res http.ResponseWriter, req *http.Request) {
	if req.Method != http.MethodPost {
		http.Error(res, "Only POST requests are allowed!", http.StatusMethodNotAllowed)
		return
	}

	originalUrl, err := io.ReadAll(req.Body)
	if err != nil {
		http.Error(res, "Key not found", http.StatusBadRequest)
		return
	}

	keyUrl := keyUrl()
	shortUrl := shortUrl(keyUrl)
	urls[keyUrl] = string(originalUrl)
	res.Write([]byte(shortUrl))
}

func keyUrl() string {
	const dictionary = "ABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	const keyLen = 8

	keyUrl := make([]byte, keyLen)
	for i := range keyUrl {
		keyUrl[i] = dictionary[rand.Intn(len(dictionary))]
	}
	return string(keyUrl)
}

func shortUrl(key string) string {
	return fmt.Sprintf("http://localhost:8080/%s", key)
}
