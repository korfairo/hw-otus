package internalhttp

import (
	"log"
	"net/http"
	"time"
)

type loggingResponseWriter struct {
	http.ResponseWriter
	statusCode int
}

func newLoggingResponseWriter(w http.ResponseWriter) *loggingResponseWriter {
	return &loggingResponseWriter{
		ResponseWriter: w,
		statusCode:     http.StatusOK,
	}
}

func (lrw *loggingResponseWriter) WriteHeader(code int) {
	lrw.statusCode = code
	lrw.ResponseWriter.WriteHeader(code)
}

func loggingMiddleware(next http.Handler) http.Handler { //nolint:unused
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		lrw := newLoggingResponseWriter(w)
		next.ServeHTTP(lrw, r)
		elapsed := time.Since(start)
		log.Printf("--> %s %s %s %s %d %.6f %s\n",
			r.RemoteAddr, r.Method, r.RequestURI, r.Proto, lrw.statusCode, elapsed.Seconds(), r.UserAgent())
	})
}
