package main

import (
	"net/http"
)

type ResponseWriter struct {
	w http.ResponseWriter
	status int
}

func (w *ResponseWriter) Header() http.Header {
	return w.w.Header()
}

func (w *ResponseWriter) Write(p []byte) (int, error) {
	if w.status == 0 {
		w.status = http.StatusOK
	}

	return w.w.Write(p)
}

func (w *ResponseWriter) WriteHeader(s int) {
	if w.status == 0 {
		w.status = s
	}

	w.w.WriteHeader(s)
}

func (w *ResponseWriter) Status() int {
	return w.status
}

