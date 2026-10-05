package middleware

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestMetricsMiddlewarePreservesResponse(t *testing.T) {
	handler := Metrics(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusCreated)
		_, _ = w.Write([]byte("created"))
	}))

	req := httptest.NewRequest(http.MethodGet, "/test", nil)
	rec := httptest.NewRecorder()

	handler.ServeHTTP(rec, req)

	if rec.Code != http.StatusCreated {
		t.Fatalf("expected status 201, got %d", rec.Code)
	}

	if rec.Body.String() != "created" {
		t.Fatalf("expected body created, got %s", rec.Body.String())
	}
}
