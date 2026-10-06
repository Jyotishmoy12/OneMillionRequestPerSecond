package config

import "testing"

func TestLoadUsesDefaults(t *testing.T) {
	t.Setenv("PORT", "")
	t.Setenv("DATABASE_URL", "")
	t.Setenv("REDIS_ADDR", "")

	cfg := Load()

	if cfg.HTTPAddr != ":8080" {
		t.Fatalf("expected default addr :8080, got %s", cfg.HTTPAddr)
	}

	expectedDatabaseURL := "postgres://postgres:postgres@localhost:5432/onemillionrps?sslmode=disable"
	if cfg.DatabaseURL != expectedDatabaseURL {
		t.Fatalf("expected default database url %s, got %s", expectedDatabaseURL, cfg.DatabaseURL)
	}

	if cfg.RedisAddr != "localhost:6379" {
		t.Fatalf("expected default redis addr localhost:6379, got %s", cfg.RedisAddr)
	}

	if cfg.DBMaxConns != 10 {
		t.Fatalf("expected default db max conns 10, got %d", cfg.DBMaxConns)
	}

	if cfg.DBMinConns != 2 {
		t.Fatalf("expected default db min conns 2, got %d", cfg.DBMinConns)
	}
}

func TestLoadUsesEnvironmentValues(t *testing.T) {
	t.Setenv("PORT", "9090")
	t.Setenv("DATABASE_URL", "postgres://user:pass@db:5432/app?sslmode=disable")
	t.Setenv("REDIS_ADDR", "redis:6379")
	t.Setenv("DB_MAX_CONNS", "25")
	t.Setenv("DB_MIN_CONNS", "5")

	cfg := Load()

	if cfg.HTTPAddr != ":9090" {
		t.Fatalf("expected addr :9090, got %s", cfg.HTTPAddr)
	}

	if cfg.DatabaseURL != "postgres://user:pass@db:5432/app?sslmode=disable" {
		t.Fatalf("unexpected database url %s", cfg.DatabaseURL)
	}

	if cfg.RedisAddr != "redis:6379" {
		t.Fatalf("expected redis addr redis:6379, got %s", cfg.RedisAddr)
	}

	if cfg.DBMaxConns != 25 {
		t.Fatalf("expected db max conns 25, got %d", cfg.DBMaxConns)
	}

	if cfg.DBMinConns != 5 {
		t.Fatalf("expected db min conns 5, got %d", cfg.DBMinConns)
	}
}

func TestLoadCapsMinConnsToMaxConns(t *testing.T) {
	t.Setenv("DB_MAX_CONNS", "5")
	t.Setenv("DB_MIN_CONNS", "20")

	cfg := Load()

	if cfg.DBMaxConns != 5 {
		t.Fatalf("expected db max conns 5, got %d", cfg.DBMaxConns)
	}

	if cfg.DBMinConns != 5 {
		t.Fatalf("expected db min conns capped to 5, got %d", cfg.DBMinConns)
	}
}
