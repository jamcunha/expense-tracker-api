package middleware

import (
	"log"
	"net/http"
	"time"
)

type wrappedWritter struct {
	http.ResponseWriter
	statusCode int
}

func (w *wrappedWritter) WriteHeader(statusCode int) {
	w.statusCode = statusCode
	w.ResponseWriter.WriteHeader(statusCode)
}

func Logging(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()

		wrapped := &wrappedWritter{w, http.StatusOK}
		next.ServeHTTP(wrapped, r)

		log.Println(r.Method, wrapped.statusCode, r.URL.Path, time.Since(start))
	})
}
