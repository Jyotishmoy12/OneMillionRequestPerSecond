package middleware

import (
	"net/http"
	"strconv"
	"strings"
	"time"

	"onemillionrps/internal/metrics"
)

type statusRecorder struct {
	http.ResponseWriter
	statusCode int
}

func newStatusRecorder(w http.ResponseWriter) *statusRecorder {
	return &statusRecorder{
		ResponseWriter: w,
		statusCode:     http.StatusOK,
	}
}

func (r *statusRecorder) WriteHeader(statusCode int) {
	r.statusCode = statusCode
	r.ResponseWriter.WriteHeader(statusCode)
}

func Metrics(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		startedAt := time.Now()

		recorder := newStatusRecorder(w)
		next.ServeHTTP(recorder, r)

		status := strconv.Itoa(recorder.statusCode)
		path := r.Pattern // Use the registered pattern as the path label

		if path == "" {
			path = r.URL.Path // Fallback to the actual URL path if no pattern is registered
		} else if _, routePath, found := strings.Cut(path, " "); found {
			path = routePath
		}

		metrics.HTTPRequestsTotal.WithLabelValues(r.Method, path, status).Inc()
		metrics.HTTPRequestDurationSeconds.WithLabelValues(r.Method, path, status).Observe(time.Since(startedAt).Seconds())
	})
}
