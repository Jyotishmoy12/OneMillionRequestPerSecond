package middleware

import (
	"bytes"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestRequestLoggerAddsRequestIDHeader(t *testing.T) {
	var buffer bytes.Buffer

	logger := slog.New(slog.NewJSONHandler(&buffer, nil))

	handler := RequestLogger(logger, http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte("ok"))
	}))

	req := httptest.NewRequest(http.MethodGet, "/test", nil)
	rec := httptest.NewRecorder()

	handler.ServeHTTP(rec, req)

	requestID := rec.Header().Get(RequestIDHeader)
	if requestID == "" {
		t.Fatal("expected request id header to be set")
	}

	if rec.Code != http.StatusOK {
		t.Fatalf("expected status 200, got %d", rec.Code)
	}

	logLine := buffer.String()
	if !strings.Contains(logLine, "http request completed") {
		t.Fatalf("expected request log, got %s", logLine)
	}
}
