package main

import (
	"context"
	"log"
	"log/slog"
	"net/http"
	"os"
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

	db, err := database.Connect(ctx, cfg.DatabaseURL)
	if err != nil {
		log.Fatalf("database connection failed: %v", err)
	}
	defer db.Close()

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

	log.Printf("api server listening on %s", cfg.HTTPAddr)

	if err := server.ListenAndServe(); err != nil {
		log.Fatal(err)
	}
}
