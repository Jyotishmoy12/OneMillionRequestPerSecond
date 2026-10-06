package middleware

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestRoutePatternRemovesMethodPrefix(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "/v1/items/1", nil)
	req.Pattern = "GET /v1/items/{id}"

	actual := routePattern(req)

	if actual != "/v1/items/{id}" {
		t.Fatalf("expected /v1/items/{id}, got %s", actual)
	}
}

func TestRoutePatternFallsBackToURLPath(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "/unknown", nil)

	actual := routePattern(req)

	if actual != "/unknown" {
		t.Fatalf("expected /unknown, got %s", actual)
	}
}
