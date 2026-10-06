package metrics

import (
	"context"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/prometheus/client_golang/prometheus"
)

var DatabasePoolAcquiredConns = prometheus.NewGauge(
	prometheus.GaugeOpts{
		Name: "database_pool_acquired_conns",
		Help: "Number of database connections currently acquired from the pool.",
	},
)

var DatabasePoolIdleConns = prometheus.NewGauge(
	prometheus.GaugeOpts{
		Name: "database_pool_idle_conns",
		Help: "Number of idle database connections currently in the pool.",
	},
)

var DatabasePoolTotalConns = prometheus.NewGauge(
	prometheus.GaugeOpts{
		Name: "database_pool_total_conns",
		Help: "Total number of database connections currently in the pool.",
	},
)

var DatabasePoolWaitCount = prometheus.NewGauge(
	prometheus.GaugeOpts{
		Name: "database_pool_wait_count",
		Help: "Total number of times a request waited for a database connection.",
	},
)

var DatabasePoolWaitDurationSeconds = prometheus.NewGauge(
	prometheus.GaugeOpts{
		Name: "database_pool_wait_duration_seconds",
		Help: "Total time spent waiting for database connections in seconds.",
	},
)

func RegisterDatabasePoolMetrics() {
	prometheus.MustRegister(DatabasePoolAcquiredConns)
	prometheus.MustRegister(DatabasePoolIdleConns)
	prometheus.MustRegister(DatabasePoolTotalConns)
	prometheus.MustRegister(DatabasePoolWaitCount)
	prometheus.MustRegister(DatabasePoolWaitDurationSeconds)
}

func ObserveDatabasePool(ctx context.Context, pool *pgxpool.Pool, interval time.Duration) {
	ticker := time.NewTicker(interval)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			return

		case <-ticker.C:
			stat := pool.Stat()

			DatabasePoolAcquiredConns.Set(float64(stat.AcquiredConns()))
			DatabasePoolIdleConns.Set(float64(stat.IdleConns()))
			DatabasePoolTotalConns.Set(float64(stat.TotalConns()))
			DatabasePoolWaitCount.Set(float64(stat.EmptyAcquireCount()))
			DatabasePoolWaitDurationSeconds.Set(stat.EmptyAcquireWaitTime().Seconds())
		}
	}
}
