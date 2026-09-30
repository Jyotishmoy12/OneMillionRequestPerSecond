package cache

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/redis/go-redis/v9"

	"onemillionrps/internal/model"
)

const itemCacheTTL = 5 * time.Minute

type RedisItemCache struct {
	client *redis.Client
}

func NewRedisItemCache(client *redis.Client) *RedisItemCache {
	return &RedisItemCache{
		client: client,
	}
}

func (c *RedisItemCache) Get(ctx context.Context, id int64) (model.Item, bool, error) {
	value, err := c.client.Get(ctx, itemCacheKey(id)).Result()

	if err == redis.Nil {
		return model.Item{}, false, nil
	}

	if err != nil {
		return model.Item{}, false, err
	}

	var item model.Item
	if err := json.Unmarshal([]byte(value), &item); err != nil {
		return model.Item{}, false, err
	}
	return item, true, nil
}

func (c *RedisItemCache) Set(ctx context.Context, item model.Item) error {
	value, err := json.Marshal(item)
	if err != nil {
		return err
	}

	return c.client.Set(ctx, itemCacheKey(item.ID), value, itemCacheTTL).Err()
}

func itemCacheKey(id int64) string {
	return fmt.Sprintf("items:%d", id)
}
