package database

import (
	"context"
	"fmt"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
)

func Connect(ctx context.Context, databaseURL string, maxConns int32, minConns int32) (*pgxpool.Pool, error) {
	poolConfig, err := pgxpool.ParseConfig(databaseURL)

	if err != nil {
		return nil, err
	}

	poolConfig.MaxConns = maxConns
	poolConfig.MinConns = minConns
	poolConfig.MaxConnLifetime = 30 * time.Minute
	poolConfig.MaxConnIdleTime = 5 * time.Minute

	pool, err := pgxpool.NewWithConfig(ctx, poolConfig)

	if err != nil {
		return nil, err
	}

	if err := pool.Ping(ctx); err != nil {
		pool.Close()
		return nil, err
	}
	return pool, nil
}

func ConnectWithRetry(ctx context.Context, databaseURL string, maxConns int32, minConns int32, retryInterval time.Duration) (*pgxpool.Pool, error) {
	var lastErr error

	for {
		pool, err := Connect(ctx, databaseURL, maxConns, minConns)
		if err == nil {
			return pool, nil
		}

		lastErr = err

		select {
		case <-ctx.Done():
			return nil, fmt.Errorf("database connection retry timeout: %w", lastErr)
		case <-time.After(retryInterval):
		}
	}
}

func RunMigrations(ctx context.Context, pool *pgxpool.Pool) error {
	const query = `
CREATE TABLE IF NOT EXISTS items (
    id BIGSERIAL PRIMARY KEY,
    name TEXT NOT NULL,
    description TEXT NOT NULL,
    price_cents INTEGER NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX IF NOT EXISTS idx_items_created_at ON items (created_at DESC);

INSERT INTO items (id, name, description, price_cents)
VALUES
    (1, 'Mechanical Keyboard', 'Low-latency keyboard for serious typing and gaming.', 8999),
    (2, 'Gaming Mouse', 'Lightweight mouse with accurate sensor tracking.', 4999),
    (3, 'USB-C Hub', 'Compact hub with HDMI, USB, and Ethernet ports.', 3999),
    (4, 'Noise Cancelling Headphones', 'Wireless headphones with long battery life.', 12999),
    (5, '4K Monitor', 'High-resolution display for development and content creation.', 27999)
ON CONFLICT (id) DO NOTHING;

SELECT setval('items_id_seq', GREATEST((SELECT MAX(id) FROM items), 1));
`

	_, err := pool.Exec(ctx, query)
	return err
}
