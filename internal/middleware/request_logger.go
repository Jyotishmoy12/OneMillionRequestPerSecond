package middleware

import (
	"log/slog"
	"net/http"
	"strconv"
	"time"
)

const RequestIDHeader = "X-Request-ID"

func RequestLogger(logger *slog.Logger, next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		startedAt := time.Now()
		requestID := r.Header.Get(RequestIDHeader)
		if requestID == "" {
			requestID = strconv.FormatInt(startedAt.UnixNano(), 36)
		}

		recorder := newStatusRecorder(w)
		recorder.Header().Set(RequestIDHeader, requestID)

		next.ServeHTTP(recorder, r)

		logger.Info(
			"http request completed",
			"request_id", requestID,
			"method", r.Method,
			"path", r.URL.Path,
			"route", routePattern(r),
			"status", recorder.statusCode,
			"duration_ms", time.Since(startedAt).Milliseconds(),
		)
	})
}
