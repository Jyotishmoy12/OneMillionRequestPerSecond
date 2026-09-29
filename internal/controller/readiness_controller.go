package controller

import (
	"context"
	"net/http"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
)

type ReadinessController struct {
	db *pgxpool.Pool
}

func NewReadinessController(db *pgxpool.Pool) *ReadinessController {
	return &ReadinessController{
		db: db,
	}
}

func (c *ReadinessController) RegisterRoutes(mux *http.ServeMux) {
	mux.HandleFunc("GET /readyz", c.Ready)
}

func (c *ReadinessController) Ready(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(r.Context(), time.Second)

	defer cancel()

	if err := c.db.Ping(ctx); err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusServiceUnavailable)
		_, _ = w.Write([]byte(`{"status":"not_ready"}`))
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	_, _ = w.Write([]byte(`{"status":"ready"}`))
}
