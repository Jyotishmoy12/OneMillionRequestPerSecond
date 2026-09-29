CREATE TABLE IF NOT EXISTS items (
    id BIGSERIAL PRIMARY KEY,
    name TEXT NOT NULL,
    description TEXT NOT NULL,
    price_cents INTEGER NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX IF NOT EXISTS idx_items_created_at ON items (created_at DESC);

INSERT INTO items (name, description, price_cents)
VALUES
    ('Mechanical Keyboard', 'Low-latency keyboard for serious typing and gaming.', 8999),
    ('Gaming Mouse', 'Lightweight mouse with accurate sensor tracking.', 4999),
    ('USB-C Hub', 'Compact hub with HDMI, USB, and Ethernet ports.', 3999),
    ('Noise Cancelling Headphones', 'Wireless headphones with long battery life.', 12999),
    ('4K Monitor', 'High-resolution display for development and content creation.', 27999)
ON CONFLICT DO NOTHING;