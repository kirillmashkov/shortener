package middleware

import (
	"net/http"
)

type ResponseWrapper interface {
	http.ResponseWriter
	Status() int
	Bytes() int
}

type writer struct {
    http.ResponseWriter
    code int
    bytes int
}

func (wr *writer) Write(buf []byte) (n int, err error) {
	n, err = wr.ResponseWriter.Write(buf)
	wr.bytes += n
	return n, err
}

func (wr *writer) WriteHeader(code int) {
	wr.code = code
	wr.ResponseWriter.WriteHeader(code)
}

func Wrap(w http.ResponseWriter) ResponseWrapper {
	wrapper := writer{ResponseWriter: w}
    return &wrapper
}

func (wr *writer) Status() int {
	return wr.code
}

func (wr *writer) Bytes() int {
	return wr.bytes
}
