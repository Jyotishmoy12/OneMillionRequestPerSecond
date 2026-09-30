package service

import (
	"context"
	"log"

	"onemillionrps/internal/cache"
	"onemillionrps/internal/model"
	"onemillionrps/internal/repository"
)

type ItemService struct {
	repository repository.ItemRepository
	cache      cache.ItemCache
}

func NewItemService(repository repository.ItemRepository, itemCache cache.ItemCache) *ItemService {
	return &ItemService{
		repository: repository,
		cache:      itemCache,
	}
}

func (s *ItemService) GetByID(ctx context.Context, id int64) (model.Item, error) {
	if s.cache != nil {
		item, found, err := s.cache.Get(ctx, id)
		if err == nil && found {
			return item, nil
		}

		if err != nil {
			log.Printf("item cache get failed: %v", err)
		}
	}

	item, err := s.repository.FindByID(ctx, id)
	if err != nil {
		return model.Item{}, err
	}

	if s.cache != nil {
		if err := s.cache.Set(ctx, item); err != nil {
			log.Printf("item cache set failed: %v", err)
		}
	}

	return item, nil
}
