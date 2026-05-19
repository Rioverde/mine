package repo

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"

	"github.com/Rioverde/mine/internal/domain"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v4"
	"github.com/jackc/pgx/v4/pgxpool"
)

type ItemRepository interface {
	SaveItem(ctx context.Context, id uuid.UUID, factoryName string, it *domain.Item, payload map[string]any) error
}

type pgItemRepo struct {
	pool *pgxpool.Pool
}

func (r *pgItemRepo) SaveItem(ctx context.Context, id uuid.UUID, factoryName string, it *domain.Item, payload map[string]any) error {
	payloadJSON, err := json.Marshal(payload)
	if err != nil {
		return fmt.Errorf("marshal payload: %w", err)
	}

	tx, err := r.pool.BeginTx(ctx, pgx.TxOptions{})
	if err != nil {
		return err
	}
	defer tx.Rollback(ctx)

	if _, err := tx.Exec(ctx, `
        INSERT INTO items (id, name, quality, ore_material, ore_capacity, ingot_quality, factory)
        VALUES ($1, $2, $3, $4, $5, $6, $7)
    `,
		id,
		it.Name(), it.Quality,
		strings.ToLower(it.Ingot.Ore.Name()),
		it.Ingot.Ore.Capacity,
		it.Ingot.Quality,
		factoryName,
	); err != nil {
		return fmt.Errorf("insert items: %w", err)
	}

	if _, err := tx.Exec(ctx, `
        INSERT INTO outbox (item_id, payload) VALUES ($1, $2)
    `, id, payloadJSON); err != nil {
		return fmt.Errorf("insert outbox: %w", err)
	}

	return tx.Commit(ctx)
}
