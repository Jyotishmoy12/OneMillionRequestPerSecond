package main

import (
	"context"
	"log"
	"net/http"
	"time"

	"onemillionrps/internal/config"
	"onemillionrps/internal/controller"
	"onemillionrps/internal/database"
	"onemillionrps/internal/repository"
	"onemillionrps/internal/service"
)

func main() {
	cfg := config.Load()

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	db, err := database.Connect(ctx, cfg.DatabaseURL)
	if err != nil {
		log.Fatalf("database connection failed: %v", err)
	}
	defer db.Close()

	mux := http.NewServeMux()

	healthController := controller.NewHealthController()
	healthController.RegisterRoutes(mux)

	readinessController := controller.NewReadinessController(db)
	readinessController.RegisterRoutes(mux)

	itemRepository := repository.NewPostgresItemRepository(db)
	itemService := service.NewItemService(itemRepository)
	itemController := controller.NewItemController(itemService)

	itemController.RegisterRoutes(mux)

	server := &http.Server{
		Addr:              cfg.HTTPAddr,
		Handler:           mux,
		ReadHeaderTimeout: 2 * time.Second,
	}

	log.Printf("api server listening on %s", cfg.HTTPAddr)

	if err := server.ListenAndServe(); err != nil {
		log.Fatal(err)
	}
}
