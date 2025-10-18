package util

import (
	"bytes"
	"net/http"
)

type ResponseWriter struct {
	http.ResponseWriter
	Body        *bytes.Buffer
	WroteHeader bool
}

func (rw *ResponseWriter) WriteHeader(statusCode int) {
	if rw.WroteHeader {
		return
	}
	rw.ResponseWriter.WriteHeader(statusCode)
	rw.WroteHeader = true
}

func (rw *ResponseWriter) Write(b []byte) (int, error) {
	return rw.Body.Write(b)
}