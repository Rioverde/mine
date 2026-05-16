package storage

import (
	"context"
	"fmt"
	"os"
	"strings"
	"time"

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

func FetchNextBatch(
	ctx context.Context,
	pool *pgxpool.Pool,
	limit int,
) ([]*ItemDTO, error) {

	var lastTime time.Time
	var lastID string

	err := pool.QueryRow(ctx, `
    SELECT last_timestamp, last_id 
    FROM exporter_checkpoints 
    WHERE checkpoint_name = 'auction_merchant';
	`).Scan(&lastTime, &lastID)

	if err != nil {
		return nil, err
	}

	rows, err := pool.Query(ctx, `
        SELECT id, name, quality, ore_material, factory, created_at 
        FROM items 
        WHERE (created_at, id) > ($1, $2)
        ORDER BY created_at ASC, id ASC 
        LIMIT $3;
    `, lastTime, lastID, limit)

	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var totalRows int

	aggrErr := pool.QueryRow(ctx, `
    SELECT COUNT(*) 
    FROM items 
    WHERE (created_at, id) > ($1, $2);
`, lastTime, lastID).Scan(&totalRows)

	if aggrErr != nil {
		return nil, aggrErr
	}

	items := make([]*ItemDTO, 0, limit)

	for rows.Next() {
		var id, name, material, factory string
		var quality int
		var createdAt time.Time

		if err := rows.Scan(&id, &name, &quality, &material, &factory, &createdAt); err != nil {
			return nil, err
		}

		it := &ItemDTO{
			id:        id,
			name:      name,
			material:  material,
			factory:   factory,
			quality:   quality,
			createdAt: createdAt,
		}

		items = append(items, it)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	if len(items) != 0 {

		lastItem := items[len(items)-1]

		_, err := pool.Exec(ctx, `
            UPDATE exporter_checkpoints 
            SET
                last_timestamp = $1,
                last_id = $2,
                updated_at = now()
            WHERE checkpoint_name = 'auction_merchant';
        `, lastItem.createdAt, lastItem.id)

		if err != nil {
			return nil, err
		}
	}

	return items, nil
}

func env(key, fallback string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return fallback
}
