package middleware

import (
	"compress/gzip"
	"io"
	"net/http"
	"strings"
)

func GzipDecompress(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Header.Get("Content-Encoding") == "gzip" {
			gz, err := gzip.NewReader(r.Body)
			if err != nil {
				http.Error(w, "failed to decompress request", http.StatusBadRequest)
				return
			}
			defer gz.Close()

			r.Body = io.NopCloser(gz)
		}

		next.ServeHTTP(w, r)
	})
}

func GzipCompress(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if !strings.Contains(r.Header.Get("Accept-Encoding"), "gzip") {
			next.ServeHTTP(w, r)
			return
		}

		gzw := &gzipResponseWriter{
			ResponseWriter: w,
			writer:         nil,
		}
		defer gzw.Close()

		next.ServeHTTP(gzw, r)
	})
}

type gzipResponseWriter struct {
	http.ResponseWriter
	writer     *gzip.Writer
	headerSent bool
	statusCode int
}

func (gzw *gzipResponseWriter) Write(b []byte) (int, error) {
	contentType := gzw.Header().Get("Content-Type")

	if strings.HasPrefix(contentType, "application/json") || strings.HasPrefix(contentType, "text/html") {
		if gzw.writer == nil {
			gzw.writer = gzip.NewWriter(gzw.ResponseWriter)
			gzw.Header().Set("Content-Encoding", "gzip")
			gzw.Header().Del("Content-Length")

			if !gzw.headerSent {
				if gzw.statusCode == 0 {
					gzw.statusCode = http.StatusOK
				}
				gzw.ResponseWriter.WriteHeader(gzw.statusCode)
				gzw.headerSent = true
			}
		}
		return gzw.writer.Write(b)
	}

	if !gzw.headerSent {
		if gzw.statusCode == 0 {
			gzw.statusCode = http.StatusOK
		}
		gzw.ResponseWriter.WriteHeader(gzw.statusCode)
		gzw.headerSent = true
	}
	return gzw.ResponseWriter.Write(b)
}

func (gzw *gzipResponseWriter) WriteHeader(code int) {
	if gzw.headerSent {
		return
	}
	gzw.statusCode = code

	contentType := gzw.Header().Get("Content-Type")

	if strings.HasPrefix(contentType, "application/json") || strings.HasPrefix(contentType, "text/html") {
		if gzw.writer == nil {
			gzw.writer = gzip.NewWriter(gzw.ResponseWriter)
			gzw.Header().Set("Content-Encoding", "gzip")
			gzw.Header().Del("Content-Length")
		}
	}

	gzw.ResponseWriter.WriteHeader(code)
	gzw.headerSent = true
}

func (gzw *gzipResponseWriter) Close() {
	if gzw.writer != nil {
		gzw.writer.Close()
	}
}
