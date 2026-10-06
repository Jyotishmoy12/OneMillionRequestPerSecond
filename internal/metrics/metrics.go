package metrics

import "github.com/prometheus/client_golang/prometheus"

var HTTPRequestsTotal = prometheus.NewCounterVec(
	prometheus.CounterOpts{
		Name: "http_requests_total",
		Help: "Total number of HTTP requests.",
	},
	[]string{"method", "path", "status"},
)

var HTTPRequestDurationSeconds = prometheus.NewHistogramVec(
	prometheus.HistogramOpts{
		Name:    "http_request_duration_seconds",
		Help:    "HTTP request duration in seconds.",
		Buckets: prometheus.DefBuckets,
	},
	[]string{"method", "path", "status"},
)

var ItemCacheRequestsTotal = prometheus.NewCounterVec(
	prometheus.CounterOpts{
		Name: "item_cache_requests_total",
		Help: "Total number of item cache lookups by result.",
	},
	[]string{"result"},
)

func Register() {
	prometheus.MustRegister(HTTPRequestsTotal)
	prometheus.MustRegister(HTTPRequestDurationSeconds)
	prometheus.MustRegister(ItemCacheRequestsTotal)
	RegisterDatabasePoolMetrics()
}
