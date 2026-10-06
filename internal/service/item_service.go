package service

import (
	"context"
	"log"

	"onemillionrps/internal/cache"
	"onemillionrps/internal/metrics"
	"onemillionrps/internal/model"
	"onemillionrps/internal/repository"
)

type ItemService struct {
	repository repository.ItemRepository
	cache      cache.ItemCache
}

type UpdateItemInput struct {
	Name        string
	Description string
	PriceCents  int
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
		if err != nil {
			metrics.ItemCacheRequestsTotal.WithLabelValues("error").Inc()
			log.Printf("item cache get failed: %v", err)
		} else if found {
			metrics.ItemCacheRequestsTotal.WithLabelValues("hit").Inc()
			return item, nil
		} else {
			metrics.ItemCacheRequestsTotal.WithLabelValues("miss").Inc()
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

func (s *ItemService) GetByIDWithoutCache(ctx context.Context, id int64) (model.Item, error) {
	return s.repository.FindByID(ctx, id)
}

func (s *ItemService) Update(ctx context.Context, id int64, input UpdateItemInput) (model.Item, error) {
	item, err := s.repository.Update(ctx, id, repository.UpdateItemInput{
		Name:        input.Name,
		Description: input.Description,
		PriceCents:  input.PriceCents,
	})
	if err != nil {
		return model.Item{}, err
	}

	if s.cache != nil {
		if err := s.cache.Delete(ctx, id); err != nil {
			log.Printf("item cache delete failed: %v", err)
		}
	}

	return item, nil
}
