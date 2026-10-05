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

type updateItemRequest struct {
	Name        string `json:"name"`
	Description string `json:"description"`
	PriceCents  int    `json:"price_cents"`
}

func NewItemController(service *service.ItemService) *ItemController {
	return &ItemController{
		service: service,
	}
}

func (c *ItemController) RegisterRoutes(mux *http.ServeMux) {
	mux.HandleFunc("GET /v1/items/{id}", c.GetByID)
	mux.HandleFunc("PUT /v1/items/{id}", c.Update)
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

func (c *ItemController) Update(w http.ResponseWriter, r *http.Request) {
	idText := strings.TrimSpace(r.PathValue("id"))

	id, err := strconv.ParseInt(idText, 10, 64)
	if err != nil || id <= 0 {
		writeJSON(w, http.StatusBadRequest, map[string]string{
			"error": "invalid item id",
		})
		return
	}

	var request updateItemRequest
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		writeJSON(w, http.StatusBadRequest, map[string]string{
			"error": "invalid json body",
		})
		return
	}

	request.Name = strings.TrimSpace(request.Name)
	request.Description = strings.TrimSpace(request.Description)

	if request.Name == "" {
		writeJSON(w, http.StatusBadRequest, map[string]string{
			"error": "name is required",
		})
		return
	}

	if request.Description == "" {
		writeJSON(w, http.StatusBadRequest, map[string]string{
			"error": "description is required",
		})
		return
	}

	if request.PriceCents <= 0 {
		writeJSON(w, http.StatusBadRequest, map[string]string{
			"error": "price_cents must be greater than zero",
		})
		return
	}

	item, err := c.service.Update(r.Context(), id, service.UpdateItemInput{
		Name:        request.Name,
		Description: request.Description,
		PriceCents:  request.PriceCents,
	})
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
