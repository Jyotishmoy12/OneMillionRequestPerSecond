package model

import "time"

type Item struct {
	ID          int64     `json:"id"`
	Name        string    `json:"name"`
	Description string    `json:"description"`
	PriceCents  int       `json:"price_cents"`
	CreatedAt   time.Time `json:"created_at"`
}
