package controller

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestHealthReturnsOK(t *testing.T) {
	controller := NewHealthController()

	req := httptest.NewRequest(http.MethodGet, "/healthz", nil)
	rec := httptest.NewRecorder()

	controller.Health(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("expected status 200, got %d", rec.Code)
	}

	expectedBody := `{"status":"ok"}`
	if rec.Body.String() != expectedBody {
		t.Fatalf("expected body %s, got %s", expectedBody, rec.Body.String())
	}

	contentType := rec.Header().Get("Content-Type")
	if contentType != "application/json" {
		t.Fatalf("expected application/json, got %s", contentType)
	}
}
