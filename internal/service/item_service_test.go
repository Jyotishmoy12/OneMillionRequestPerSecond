package service

import (
	"context"
	"testing"
	"time"

	"onemillionrps/internal/model"
)

type fakeItemRepository struct {
	item model.Item
	err  error
}

func (r fakeItemRepository) FindByID(ctx context.Context, id int64) (model.Item, error) {
	return r.item, r.err
}

func TestItemServiceGetByIDReturnsItem(t *testing.T) {
	expected := model.Item{
		ID:          1,
		Name:        "Mechanical Keyboard",
		Description: "Low-latency keyboard for serious typing and gaming.",
		PriceCents:  8999,
		CreatedAt:   time.Now(),
	}

	service := NewItemService(fakeItemRepository{
		item: expected,
	})

	actual, err := service.GetByID(context.Background(), 1)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if actual.ID != expected.ID {
		t.Fatalf("expected id %d, got %d", expected.ID, actual.ID)
	}
}
