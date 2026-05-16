package storage

import (
	"context"
	"fmt"
	"os"
	"strings"

	"github.com/Rioverde/mine/internal/item"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v4/pgxpool"
)

func Connect(ctx context.Context) (*pgxpool.Pool, error) {
	dsn := fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=disable",
		env("STORAGE_DB_USER", "mine"),
		os.Getenv("STORAGE_DB_PASSWORD"),
		env("STORAGE_DB_HOST", "localhost"),
		env("STORAGE_DB_PORT", "5432"),
		env("STORAGE_DB_NAME", "mine"),
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

func env(key, fallback string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return fallback
}
