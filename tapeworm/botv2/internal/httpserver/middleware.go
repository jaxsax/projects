package httpserver

import (
	"net/http"

	"github.com/felixge/httpsnoop"
	"github.com/go-logr/logr"
	"github.com/segmentio/ksuid"
)

func loggingMiddleware(logger logr.Logger) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			requestID := r.Header.Get("x-request-id")
			if requestID == "" {
				requestID = ksuid.New().String()
			}

			lg := logger.WithValues("request_id", requestID)
			ctx := logr.NewContext(r.Context(), lg)
			w.Header().Set("x-request-id", requestID)

			newRequest := r.WithContext(ctx)
			m := httpsnoop.CaptureMetrics(next, w, newRequest)
			lg.Info(
				"request completed",
				"status_code", m.Code,
				"duration", m.Duration,
				"duration_ms", m.Duration.Milliseconds(),
				"bytes_written", m.Written,
				"method", newRequest.Method,
				"url", newRequest.URL.String(),
			)
		})
	}
}
