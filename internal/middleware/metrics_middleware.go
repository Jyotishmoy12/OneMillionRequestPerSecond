package middleware

import (
	"net/http"
	"onemillionrps/internal/metrics"
	"strconv"
	"time"
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
		path := r.URL.Path

		metrics.HTTPRequestsTotal.WithLabelValues(r.Method, path, status).Inc()
		metrics.HTTPRequestDurationSeconds.WithLabelValues(r.Method, path, status).Observe(time.Since(startedAt).Seconds())
	})
}
