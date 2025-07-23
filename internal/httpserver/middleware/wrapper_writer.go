package middleware

import (
	"net/http"
)

// Writer - тип, обертка над стандартным writer
type Writer struct {
	http.ResponseWriter
	code  int
	bytes int
}

// Writer - запись ответа
func (wr *Writer) Write(buf []byte) (n int, err error) {
	n, err = wr.ResponseWriter.Write(buf)
	wr.bytes += n
	return n, err
}

// WriteHeader - запись заголовков
func (wr *Writer) WriteHeader(code int) {
	wr.code = code
	wr.ResponseWriter.WriteHeader(code)
}

// Status - сохранение статуса в отдельное поле
func (wr *Writer) Status() int {
	return wr.code
}

// Bytes - сохранение кол-ва записанных байт в отдельное поле
func (wr *Writer) Bytes() int {
	return wr.bytes
}
