package repository

import (
	"context"

	"onemillionrps/internal/model"
)

type UpdateItemInput struct {
	Name        string
	Description string
	PriceCents  int
}

type ItemRepository interface {
	FindByID(ctx context.Context, id int64) (model.Item, error)
	Update(ctx context.Context, id int64, input UpdateItemInput) (model.Item, error)
}
