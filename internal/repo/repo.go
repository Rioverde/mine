package repo

import (
	"context"
	"fmt"
	"os"

	"github.com/jackc/pgx/v4/pgxpool"
)

// Connection owns the database pool and hands out scoped repositories.
type Connection struct {
	pool *pgxpool.Pool
}

func Connect(ctx context.Context) (*Connection, error) {
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
	return &Connection{pool: pool}, nil
}

func (c *Connection) Close() {
	c.pool.Close()
}

// Items returns a narrow repository used by the producer pipeline.
func (c *Connection) Items() ItemRepository {
	return &pgItemRepo{pool: c.pool}
}

// Outbox returns a narrow repository used by the publisher and cleaner.
func (c *Connection) Outbox() OutboxRepository {
	return &pgOutboxRepo{pool: c.pool}
}

func env(key, fallback string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return fallback
}
