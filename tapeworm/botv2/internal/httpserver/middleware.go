package httpserver

import (
	"net/http"

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

			newRequest := r.WithContext(ctx)
			next.ServeHTTP(w, newRequest)
		})
	}
}
