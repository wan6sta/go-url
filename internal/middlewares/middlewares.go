package middlewares

import (
	"fmt"
	"log/slog"
	"net/http"
	"time"
)

type HTTPMiddlewares struct {
	log *slog.Logger
}

type loggingResponseWriter struct {
	status int
	size   int
	http.ResponseWriter
}

func (r *loggingResponseWriter) Write(b []byte) (int, error) {
	// записываем ответ, используя оригинальный http.ResponseWriter
	size, err := r.ResponseWriter.Write(b)
	r.size += size // захватываем размер
	return size, err
}

func (r *loggingResponseWriter) WriteHeader(statusCode int) {
	// записываем код статуса, используя оригинальный http.ResponseWriter
	r.ResponseWriter.WriteHeader(statusCode)
	r.status = statusCode // захватываем код статуса
}

func NewHTTPMiddlewares(log *slog.Logger) *HTTPMiddlewares {
	return &HTTPMiddlewares{log: log}
}

func (hm *HTTPMiddlewares) Log(h http.Handler) http.Handler {
	logFn := func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()

		lw := &loggingResponseWriter{
			status:         0,
			size:           0,
			ResponseWriter: w,
		}

		h.ServeHTTP(lw, r)

		duration := fmt.Sprintf("%fs", time.Since(start).Seconds())
		size := fmt.Sprintf("%db", lw.size)

		hm.log.Info(
			"req info",
			"uri", r.RequestURI,
			"method", r.Method,
			"status", lw.status,
			"duration", duration,
			"size", size,
		)
	}

	return http.HandlerFunc(logFn)
}
