package middlewares

import (
	"bytes"
	"log"
	"net/http"
	"time"
)

type LogResponseWriter struct {
	http.ResponseWriter
	statusCode int
	buf        bytes.Buffer
}

func NewLogResponseWriter(w http.ResponseWriter) *LogResponseWriter {
	return &LogResponseWriter{ResponseWriter: w}
}

func (w *LogResponseWriter) WriteHeader(code int) {
	w.statusCode = code
	w.ResponseWriter.WriteHeader(code)
}

func (w *LogResponseWriter) Write(body []byte) (int, error) {
	w.buf.Write(body)
	return w.ResponseWriter.Write(body)
}

// Log middleware log each received request an given response.
func Log(handler http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		startTime := time.Now()

		logRw := NewLogResponseWriter(w)
		handler.ServeHTTP(logRw, r)

		log.Printf("%s %s duration=%s status=%d", r.Method, r.RequestURI, time.Since(startTime).String(), logRw.statusCode)
	})
}
