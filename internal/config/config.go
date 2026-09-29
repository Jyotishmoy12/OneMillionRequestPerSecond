package config

import (
	"os"
	"strings"
)

type Config struct {
	HTTPAddr    string
	DatabaseURL string
}

func Load() Config {
	port := strings.TrimSpace(os.Getenv("PORT"))
	if port == "" {
		port = "8080"
	}

	databaseURL := strings.TrimSpace(os.Getenv("DATABASE_URL"))
	if databaseURL == "" {
		databaseURL = "postgres://postgres:postgres@localhost:5432/onemillionrps?sslmode=disable"
	}

	return Config{
		HTTPAddr:    ":" + port,
		DatabaseURL: databaseURL,
	}
}
