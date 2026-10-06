package main

import (
	"context"
	"errors"
	"log"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/redis/go-redis/v9"

	"onemillionrps/internal/cache"
	"onemillionrps/internal/config"
	"onemillionrps/internal/controller"
	"onemillionrps/internal/database"
	"onemillionrps/internal/metrics"
	"onemillionrps/internal/middleware"
	"onemillionrps/internal/repository"
	"onemillionrps/internal/service"
)

func main() {

	logger := slog.New(slog.NewJSONHandler(os.Stdout, nil))

	slog.SetDefault(logger)
	cfg := config.Load()

	metrics.Register()

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	db, err := database.Connect(ctx, cfg.DatabaseURL, cfg.DBMaxConns, cfg.DBMinConns)
	if err != nil {
		log.Fatalf("database connection failed: %v", err)
	}
	defer db.Close()
	logger.Info(
		"database pool configured",
		"max_conns", cfg.DBMaxConns,
		"min_conns", cfg.DBMinConns,
	)

	metricsCtx, cancelMetrics := context.WithCancel(context.Background())
	defer cancelMetrics()

	go metrics.ObserveDatabasePool(metricsCtx, db, 10*time.Second)

	redisClient := redis.NewClient(&redis.Options{
		Addr: cfg.RedisAddr,
	})
	defer redisClient.Close()

	mux := http.NewServeMux()

	healthController := controller.NewHealthController()
	healthController.RegisterRoutes(mux)

	readinessController := controller.NewReadinessController(db)
	readinessController.RegisterRoutes(mux)

	itemRepository := repository.NewPostgresItemRepository(db)
	itemCache := cache.NewRedisItemCache(redisClient)
	itemService := service.NewItemService(itemRepository, itemCache)
	itemController := controller.NewItemController(itemService)

	itemController.RegisterRoutes(mux)

	mux.Handle("GET /metrics", promhttp.Handler())

	server := &http.Server{
		Addr:              cfg.HTTPAddr,
		Handler:           middleware.RequestLogger(logger, middleware.Metrics(mux)),
		ReadHeaderTimeout: 2 * time.Second,
	}

	serverErrors := make(chan error, 1)

	go func() {
		logger.Info("api server listening", "addr", cfg.HTTPAddr)

		if err := server.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			serverErrors <- err
		}
	}()

	shutdownCtx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)

	defer stop()

	select {
	case err := <-serverErrors:
		logger.Error("api server error", "error", err)
	case <-shutdownCtx.Done():
		logger.Info("shutdown signal received")
	}

	timeoutCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := server.Shutdown(timeoutCtx); err != nil {
		logger.Error("api server shutdown error", "error", err)
	} else {
		logger.Info("api server shutdown gracefully")
	}
}
