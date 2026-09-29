package controller

import (
	"encoding/json"
	"errors"
	"net/http"
	"strconv"
	"strings"

	"github.com/jackc/pgx/v5"

	"onemillionrps/internal/service"
)

type ItemController struct {
	service *service.ItemService
}

func NewItemController(service *service.ItemService) *ItemController {
	return &ItemController{
		service: service,
	}
}

func (c *ItemController) RegisterRoutes(mux *http.ServeMux) {
	mux.HandleFunc("GET /v1/items/{id}", c.GetByID)
}

func (c *ItemController) GetByID(w http.ResponseWriter, r *http.Request) {
	idText := strings.TrimSpace(r.PathValue("id"))

	id, err := strconv.ParseInt(idText, 10, 64)
	if err != nil || id <= 0 {
		writeJSON(w, http.StatusBadRequest, map[string]string{
			"error": "invalid item id",
		})
		return
	}

	item, err := c.service.GetByID(r.Context(), id)
	if errors.Is(err, pgx.ErrNoRows) {
		writeJSON(w, http.StatusNotFound, map[string]string{
			"error": "item not found",
		})
		return
	}

	if err != nil {
		writeJSON(w, http.StatusInternalServerError, map[string]string{
			"error": "internal server error",
		})
		return
	}

	writeJSON(w, http.StatusOK, item)
}

func writeJSON(w http.ResponseWriter, statusCode int, value any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)

	_ = json.NewEncoder(w).Encode(value)
}
