package service

import (
	"context"

	"onemillionrps/internal/model"
	"onemillionrps/internal/repository"
)

type ItemService struct {
	repository repository.ItemRepository
}

func NewItemService(repository repository.ItemRepository) *ItemService {
	return &ItemService{
		repository: repository,
	}
}

func (s *ItemService) GetByID(ctx context.Context, id int64) (model.Item, error) {
	return s.repository.FindByID(ctx, id)
}
