package repo

import (
	"context"
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/Rioverde/mine/internal/domain"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v4/pgxpool"
)

type Repo interface {
	Save(ctx context.Context, factoryName string, it *domain.Item) (uuid.UUID, error)
	FetchNextBatch(ctx context.Context, limit int) ([]*ItemDTO, error)
	Close()
}

type pgxRepo struct {
	pool *pgxpool.Pool
}

func Connect(ctx context.Context) (Repo, error) {
	dsn := fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=disable",
		env("STORAGE_DB_USER", "mine"),
		os.Getenv("STORAGE_DB_PASSWORD"),
		env("STORAGE_DB_HOST", "localhost"),
		env("STORAGE_DB_PORT", "5432"),
		env("STORAGE_DB_NAME", "mine"),
	)
	pool, err := pgxpool.Connect(ctx, dsn)
	if err != nil {
		return nil, err
	}
	return &pgxRepo{pool: pool}, nil
}

func (r *pgxRepo) Close() {
	r.pool.Close()
}

func (r *pgxRepo) Save(ctx context.Context, factoryName string, it *domain.Item) (uuid.UUID, error) {
	id := uuid.New()
	_, err := r.pool.Exec(ctx, `
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

func (r *pgxRepo) FetchNextBatch(ctx context.Context, limit int) ([]*ItemDTO, error) {
	var lastTime time.Time
	var lastID string

	err := r.pool.QueryRow(ctx, `
        SELECT last_timestamp, last_id
        FROM exporter_checkpoints
        WHERE checkpoint_name = 'auction_merchant';
    `).Scan(&lastTime, &lastID)
	if err != nil {
		return nil, err
	}

	rows, err := r.pool.Query(ctx, `
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

	items := make([]*ItemDTO, 0, limit)
	for rows.Next() {
		var id, name, material, factory string
		var quality int
		var createdAt time.Time

		if err := rows.Scan(&id, &name, &quality, &material, &factory, &createdAt); err != nil {
			return nil, err
		}

		items = append(items, &ItemDTO{
			id:        id,
			name:      name,
			material:  material,
			factory:   factory,
			quality:   quality,
			createdAt: createdAt,
		})
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	if len(items) != 0 {
		last := items[len(items)-1]
		if _, err := r.pool.Exec(ctx, `
            UPDATE exporter_checkpoints
            SET
                last_timestamp = $1,
                last_id = $2,
                updated_at = now()
            WHERE checkpoint_name = 'auction_merchant';
        `, last.createdAt, last.id); err != nil {
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
