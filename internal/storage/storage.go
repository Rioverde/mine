package storage

import (
	"context"
	"fmt"
	"strings"

	"github.com/Rioverde/mine/internal/config"
	"github.com/Rioverde/mine/internal/item"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v4/pgxpool"
)

func Connect(ctx context.Context, cfg config.StorageConfig) (*pgxpool.Pool, error) {
	dsn := fmt.Sprintf("postgres://%s:%s@%s:%d/%s?sslmode=disable",
		cfg.User, cfg.Password, cfg.Host, cfg.Port, cfg.Name,
	)
	return pgxpool.Connect(ctx, dsn)
}

func SaveItem(ctx context.Context, pool *pgxpool.Pool, factoryName string, it *item.Item) (uuid.UUID, error) {
	id := uuid.New()
	_, err := pool.Exec(ctx, `
        INSERT INTO items (id, name, quality, ore_material, ore_capacity, ingot_quality, factory)
        VALUES ($1, $2, $3, $4, $5, $6, $7)`,
		id,
		it.Name(), it.Quality,
		strings.ToLower(it.Ingot.Ore.Name()),
		it.Ingot.Ore.Capacity,
		it.Ingot.Quality,
		factoryName,
	)
	return id, err
}
