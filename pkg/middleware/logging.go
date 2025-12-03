package middleware

import (
	"net/http"
	"time"

	"github.com/rs/zerolog/log"
)

func Logging(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()

		lw := &loggingResponseWriter{ResponseWriter: w, statusCode: http.StatusOK}

		next.ServeHTTP(lw, r)

		duration := time.Since(start)

		log.Info().
			Str("uri", r.RequestURI).
			Str("method", r.Method).
			Dur("duration", duration).
			Int("status", lw.statusCode).
			Int("size", lw.size).
			Msg("HTTP request")
	})
}

type loggingResponseWriter struct {
	http.ResponseWriter
	statusCode int
	size       int
}

func (lrw *loggingResponseWriter) WriteHeader(code int) {
	lrw.statusCode = code
	lrw.ResponseWriter.WriteHeader(code)
}

func (lrw *loggingResponseWriter) Write(b []byte) (int, error) {
	n, err := lrw.ResponseWriter.Write(b)
	lrw.size += n
	return n, err
}
