package config

import (
	"os"
	"strconv"
	"strings"
)

type Config struct {
	HTTPAddr    string
	DatabaseURL string
	RedisAddr   string
	DBMaxConns  int32
	DBMinConns  int32
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

	redisAddr := strings.TrimSpace(os.Getenv("REDIS_ADDR"))

	if redisAddr == "" {
		redisAddr = "localhost:6379"
	}

	dbMaxConns := readInt32Env("DB_MAX_CONNS", 10)
	dbMinConns := readInt32Env("DB_MIN_CONNS", 2)
	if dbMinConns > dbMaxConns {
		dbMinConns = dbMaxConns
	}

	return Config{
		HTTPAddr:    ":" + port,
		DatabaseURL: databaseURL,
		RedisAddr:   redisAddr,
		DBMaxConns:  dbMaxConns,
		DBMinConns:  dbMinConns,
	}
}

func readInt32Env(name string, fallback int32) int32 {
	valueStr := strings.TrimSpace(os.Getenv(name))
	if valueStr == "" {
		return fallback
	}

	parsed, err := strconv.ParseInt(valueStr, 10, 32)
	if err != nil || parsed <= 0 {
		return fallback
	}

	return int32(parsed)
}
