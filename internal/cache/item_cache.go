package cache

import (
	"context"

	"onemillionrps/internal/model"
)

type ItemCache interface {
	Get(ctx context.Context, id int64) (model.Item, bool, error)
	Set(ctx context.Context, item model.Item) error
	Delete(ctx context.Context, id int64) error
}
