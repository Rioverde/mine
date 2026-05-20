package main

import (
	"context"
	"log"
	"log/slog"
	"os/signal"
	"sync"
	"syscall"
	"time"

	"github.com/Rioverde/mine/internal/config"
	"github.com/Rioverde/mine/internal/events"
	"github.com/Rioverde/mine/internal/repo"
)

func main() {
	cfg, err := config.ReadConfig("config.yaml")
	if err != nil {
		log.Fatal(err)
	}

	ctx, cancel := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer cancel()

	conn, err := repo.Connect(ctx)
	if err != nil {
		log.Fatal(err)
	}
	defer conn.Close()

	outboxRepo := conn.NewOutboxRepository()

	bus, err := events.Connect(ctx, cfg.App.CacheStreamMaxLength)
	if err != nil {
		log.Fatal(err)
	}
	defer bus.Close()

	var wg sync.WaitGroup

	wg.Add(1)
	go func() {
		defer wg.Done()
		events.OutboxPublisher(ctx, outboxRepo, bus, &cfg.App)
	}()

	wg.Add(1)
	go func() {
		defer wg.Done()
		events.OutboxCleaner(ctx, outboxRepo, time.Hour)
	}()

	wg.Wait()
	slog.Info("Outbox service stopped.")
}
