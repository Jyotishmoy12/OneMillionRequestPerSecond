package repository

import (
	"context"

	"github.com/jackc/pgx/v5/pgxpool"

	"onemillionrps/internal/model"
)

type PostgresItemRepository struct {
	db *pgxpool.Pool
}

func NewPostgresItemRepository(db *pgxpool.Pool) *PostgresItemRepository {
	return &PostgresItemRepository{
		db: db,
	}
}

func (r *PostgresItemRepository) FindByID(ctx context.Context, id int64) (model.Item, error) {
	const query = `
		SELECT id, name, description, price_cents, created_at
		FROM items
		WHERE id = $1
	`

	var item model.Item

	err := r.db.QueryRow(ctx, query, id).Scan(
		&item.ID,
		&item.Name,
		&item.Description,
		&item.PriceCents,
		&item.CreatedAt,
	)
	if err != nil {
		return model.Item{}, err
	}

	return item, nil
}
