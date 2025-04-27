package middleware

import (
	"net/http"
	"strings"
	"time"

	log "github.com/sirupsen/logrus"
)

// RequestLoggerMiddleware returns a middleware that logs HTTP requests except for the specified paths.
// It logs method, URL, HTTP version, remote address, status code, response size, and processing duration as structured logs.
func RequestLoggerMiddleware(skipPaths []string) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			path := r.URL.Path
			for _, skipPath := range skipPaths {
				if strings.HasPrefix(path, skipPath) {
					next.ServeHTTP(w, r)
					return
				}
			}

			start := time.Now()
			ww := &responseWriter{ResponseWriter: w, statusCode: http.StatusOK}
			next.ServeHTTP(ww, r)
			duration := time.Since(start)
			log.WithFields(log.Fields{
				"method":        r.Method,
				"url":           r.URL.String(),
				"http_version":  r.Proto,
				"remote_addr":   r.RemoteAddr,
				"status":        ww.statusCode,
				"response_size": ww.size,
				"duration_ms":   duration.Milliseconds(),
			}).Info("handled request")
		})
	}
}

type responseWriter struct {
	http.ResponseWriter
	statusCode int
	size       int
}

func (rw *responseWriter) WriteHeader(code int) {
	rw.statusCode = code
	rw.ResponseWriter.WriteHeader(code)
}

func (rw *responseWriter) Write(b []byte) (int, error) {
	size, err := rw.ResponseWriter.Write(b)
	rw.size += size
	return size, err
}
