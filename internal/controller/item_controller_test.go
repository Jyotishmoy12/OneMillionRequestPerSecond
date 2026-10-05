package controller

import (
	"bytes"
	"context"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"onemillionrps/internal/model"
	"onemillionrps/internal/repository"
	"onemillionrps/internal/service"
)

type fakeControllerItemRepository struct {
	item model.Item
	err  error
}

func (r fakeControllerItemRepository) FindByID(ctx context.Context, id int64) (model.Item, error) {
	return r.item, r.err
}

func (r fakeControllerItemRepository) Update(ctx context.Context, id int64, input repository.UpdateItemInput) (model.Item, error) {
	return r.item, r.err
}

func TestItemControllerReturnsItem(t *testing.T) {
	item := model.Item{
		ID:          1,
		Name:        "Mechanical Keyboard",
		Description: "Low-latency keyboard for serious typing and gaming.",
		PriceCents:  8999,
		CreatedAt:   time.Date(2026, 9, 29, 12, 0, 0, 0, time.UTC),
	}

	itemService := service.NewItemService(fakeControllerItemRepository{
		item: item,
	}, nil)
	controller := NewItemController(itemService)

	req := httptest.NewRequest(http.MethodGet, "/v1/items/1", nil)
	req.SetPathValue("id", "1")

	rec := httptest.NewRecorder()

	controller.GetByID(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("expected status 200, got %d", rec.Code)
	}

	expectedContentType := "application/json"
	if rec.Header().Get("Content-Type") != expectedContentType {
		t.Fatalf("expected content type %s, got %s", expectedContentType, rec.Header().Get("Content-Type"))
	}
}

func TestItemControllerRejectsInvalidID(t *testing.T) {
	itemService := service.NewItemService(fakeControllerItemRepository{}, nil)
	controller := NewItemController(itemService)

	req := httptest.NewRequest(http.MethodGet, "/v1/items/abc", nil)
	req.SetPathValue("id", "abc")

	rec := httptest.NewRecorder()

	controller.GetByID(rec, req)

	if rec.Code != http.StatusBadRequest {
		t.Fatalf("expected status 400, got %d", rec.Code)
	}
}

func TestItemControllerUpdatesItem(t *testing.T) {
	item := model.Item{
		ID:          1,
		Name:        "Mechanical Keyboard Pro",
		Description: "Updated low-latency keyboard.",
		PriceCents:  9999,
		CreatedAt:   time.Date(2026, 9, 29, 12, 0, 0, 0, time.UTC),
	}

	itemService := service.NewItemService(fakeControllerItemRepository{
		item: item,
	}, nil)
	controller := NewItemController(itemService)

	body := []byte(`{"name":"Mechanical Keyboard Pro","description":"Updated low-latency keyboard.","price_cents":9999}`)
	req := httptest.NewRequest(http.MethodPut, "/v1/items/1", bytes.NewReader(body))
	req.SetPathValue("id", "1")

	rec := httptest.NewRecorder()

	controller.Update(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("expected status 200, got %d", rec.Code)
	}

	expectedContentType := "application/json"
	if rec.Header().Get("Content-Type") != expectedContentType {
		t.Fatalf("expected content type %s, got %s", expectedContentType, rec.Header().Get("Content-Type"))
	}
}

func TestItemControllerRejectsInvalidUpdateBody(t *testing.T) {
	itemService := service.NewItemService(fakeControllerItemRepository{}, nil)
	controller := NewItemController(itemService)

	body := []byte(`{"name":"","description":"Updated low-latency keyboard.","price_cents":9999}`)
	req := httptest.NewRequest(http.MethodPut, "/v1/items/1", bytes.NewReader(body))
	req.SetPathValue("id", "1")

	rec := httptest.NewRecorder()

	controller.Update(rec, req)

	if rec.Code != http.StatusBadRequest {
		t.Fatalf("expected status 400, got %d", rec.Code)
	}
}
