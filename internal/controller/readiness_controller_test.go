package controller

import (
	"net/http"
	"testing"
)

func TestReadinessControllerCanRegisterRoutes(t *testing.T) {
	mux := http.NewServeMux()

	controller := NewReadinessController(nil)
	controller.RegisterRoutes(mux)
}
