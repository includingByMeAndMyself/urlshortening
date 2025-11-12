package middleware

import (
	"log"
	"net/http"
	"time"
)

// Logging — middleware для логирования HTTP-запросов
func Logging(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()

		// Оборачиваем ResponseWriter, чтобы получить статус и размер ответа
		lw := &loggingResponseWriter{ResponseWriter: w, statusCode: http.StatusOK}

		next.ServeHTTP(lw, r)

		duration := time.Since(start)
		log.Printf("%s %s %d %d %s",
			r.Method,
			r.URL.Path,
			lw.statusCode,
			lw.size,
			duration,
		)
	})
}

// loggingResponseWriter — обёртка над http.ResponseWriter для перехвата статуса и размера тела
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
