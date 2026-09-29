package config

import "testing"

func TestLoadUsesDefaults(t *testing.T) {
	t.Setenv("PORT", "")
	t.Setenv("DATABASE_URL", "")

	cfg := Load()

	if cfg.HTTPAddr != ":8080" {
		t.Fatalf("expected default addr :8080, got %s", cfg.HTTPAddr)
	}

	expectedDatabaseURL := "postgres://postgres:postgres@localhost:5432/onemillionrps?sslmode=disable"
	if cfg.DatabaseURL != expectedDatabaseURL {
		t.Fatalf("expected default database url %s, got %s", expectedDatabaseURL, cfg.DatabaseURL)
	}
}

func TestLoadUsesEnvironmentValues(t *testing.T) {
	t.Setenv("PORT", "9090")
	t.Setenv("DATABASE_URL", "postgres://user:pass@db:5432/app?sslmode=disable")

	cfg := Load()

	if cfg.HTTPAddr != ":9090" {
		t.Fatalf("expected addr :9090, got %s", cfg.HTTPAddr)
	}

	if cfg.DatabaseURL != "postgres://user:pass@db:5432/app?sslmode=disable" {
		t.Fatalf("unexpected database url %s", cfg.DatabaseURL)
	}
}