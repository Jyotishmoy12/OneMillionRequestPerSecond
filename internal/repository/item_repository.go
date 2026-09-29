package repository

import (
	"context"
	"onemillionrps/internal/model"
)

type ItemRepository interface {
	FindByID(ctx context.Context, id int64) (model.Item, error)
}
