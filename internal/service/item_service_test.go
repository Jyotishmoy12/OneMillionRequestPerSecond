package service

import (
	"context"
	"testing"
	"time"

	"onemillionrps/internal/model"
	"onemillionrps/internal/repository"
)

type fakeItemRepository struct {
	item  model.Item
	err   error
	calls int
}

func (r *fakeItemRepository) FindByID(ctx context.Context, id int64) (model.Item, error) {
	r.calls++
	return r.item, r.err
}

func (r *fakeItemRepository) Update(ctx context.Context, id int64, input repository.UpdateItemInput) (model.Item, error) {
	r.calls++
	return r.item, r.err
}

type fakeItemCache struct {
	item    model.Item
	found   bool
	getErr  error
	setErr  error
	delErr  error
	sets    int
	deletes int
}

func (c *fakeItemCache) Get(ctx context.Context, id int64) (model.Item, bool, error) {
	return c.item, c.found, c.getErr
}

func (c *fakeItemCache) Set(ctx context.Context, item model.Item) error {
	c.sets++
	return c.setErr
}

func (c *fakeItemCache) Delete(ctx context.Context, id int64) error {
	c.deletes++
	return c.delErr
}

func TestItemServiceGetByIDReturnsCachedItem(t *testing.T) {
	cachedItem := model.Item{
		ID:          1,
		Name:        "Mechanical Keyboard",
		Description: "Low-latency keyboard for serious typing and gaming.",
		PriceCents:  8999,
		CreatedAt:   time.Now(),
	}

	repository := &fakeItemRepository{}
	itemCache := &fakeItemCache{
		item:  cachedItem,
		found: true,
	}

	service := NewItemService(repository, itemCache)

	actual, err := service.GetByID(context.Background(), 1)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if actual.ID != cachedItem.ID {
		t.Fatalf("expected id %d, got %d", cachedItem.ID, actual.ID)
	}

	if repository.calls != 0 {
		t.Fatalf("expected repository not to be called on cache hit")
	}
}

func TestItemServiceGetByIDReadsRepositoryOnCacheMiss(t *testing.T) {
	expected := model.Item{
		ID:          1,
		Name:        "Mechanical Keyboard",
		Description: "Low-latency keyboard for serious typing and gaming.",
		PriceCents:  8999,
		CreatedAt:   time.Now(),
	}

	repository := &fakeItemRepository{
		item: expected,
	}
	itemCache := &fakeItemCache{
		found: false,
	}

	service := NewItemService(repository, itemCache)

	actual, err := service.GetByID(context.Background(), 1)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if actual.ID != expected.ID {
		t.Fatalf("expected id %d, got %d", expected.ID, actual.ID)
	}

	if repository.calls != 1 {
		t.Fatalf("expected repository to be called once, got %d", repository.calls)
	}

	if itemCache.sets != 1 {
		t.Fatalf("expected cache set once, got %d", itemCache.sets)
	}
}

func TestItemServiceUpdateInvalidatesCache(t *testing.T) {
	expected := model.Item{
		ID:          1,
		Name:        "Mechanical Keyboard Pro",
		Description: "Updated low-latency keyboard.",
		PriceCents:  9999,
		CreatedAt:   time.Now(),
	}

	repository := &fakeItemRepository{
		item: expected,
	}
	itemCache := &fakeItemCache{}

	service := NewItemService(repository, itemCache)

	actual, err := service.Update(context.Background(), 1, UpdateItemInput{
		Name:        expected.Name,
		Description: expected.Description,
		PriceCents:  expected.PriceCents,
	})
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if actual.ID != expected.ID {
		t.Fatalf("expected id %d, got %d", expected.ID, actual.ID)
	}

	if repository.calls != 1 {
		t.Fatalf("expected repository to be called once, got %d", repository.calls)
	}

	if itemCache.deletes != 1 {
		t.Fatalf("expected cache delete once, got %d", itemCache.deletes)
	}
}
