package middleware

import (
	"net/http"
	"time"

	"github.com/escoutdoor/study-platform/pkg/logger"
)

type statusRecorder struct {
	http.ResponseWriter
	statusCode int
}

func (r *statusRecorder) WriteHeader(statusCode int) {
	r.statusCode = statusCode
	r.ResponseWriter.WriteHeader(statusCode)
}

func Logging(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		logger.DebugKV(r.Context(), "incoming request",
			"method", r.Method,
			"path", r.URL.Path,
		)

		recorder := &statusRecorder{
			ResponseWriter: w,
			statusCode:     http.StatusOK,
		}

		next.ServeHTTP(recorder, r)

		logger.DebugKV(r.Context(), "request completed",
			"method", r.Method,
			"path", r.URL.Path,
			"status", recorder.statusCode,
			"duration_ms", time.Since(start).Milliseconds(),
		)
	})
}
