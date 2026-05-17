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
	"github.com/Rioverde/mine/internal/domain"
	"github.com/Rioverde/mine/internal/events"
	"github.com/Rioverde/mine/internal/factory"
	"github.com/Rioverde/mine/internal/mine"
	"github.com/Rioverde/mine/internal/repo"
	"github.com/google/uuid"
)

const (
	Smithy  = "Smithy"
	Smelter = "Smelter"
)

func main() {
	cfg, err := config.ReadConfig("config.yaml")
	if err != nil {
		log.Fatal(err)
	}

	ctx, cancel := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer cancel()

	r, err := repo.Connect(ctx)
	if err != nil {
		log.Fatal(err)
	}
	defer r.Close()

	bus, err := events.Connect(ctx, cfg.App.CacheStreamMaxLength)
	if err != nil {
		log.Fatal(err)
	}
	defer bus.Close()

	ores := mine.Run(ctx, cfg)
	ingots := factory.New(factory.WithName(ctx, Smelter), cfg, ores, domain.NewSmelter(cfg.Ingot))
	items := factory.New(factory.WithName(ctx, Smithy), cfg, ingots, domain.NewSmithy(cfg.Item))

	var wg sync.WaitGroup

	wg.Add(1)
	go func() {
		defer wg.Done()
		events.OutboxPublisher(ctx, r, bus, &cfg.App)
	}()

	wg.Add(1)
	go func() {
		defer wg.Done()
		events.OutboxCleaner(ctx, r, time.Hour)
	}()

	wg.Add(1)
	go func() {
		defer wg.Done()

		defer cancel()

		for item := range items {
			id := uuid.New()
			payload := events.NewItemPayload(id.String(), Smithy, item).ToMap()

			saveCtx, cancel := context.WithTimeout(ctx, cfg.App.SaveTimeout)
			err := r.SaveItem(saveCtx, id, Smithy, item, payload)
			cancel()
			if err != nil {
				slog.Error("save item", "err", err)
				continue
			}

			slog.Info("item saved",
				"id", id,
				"type", item.Name(),
				"quality", item.Quality,
				"ore", item.Ingot.Ore.Name(),
			)
		}

		slog.Info("DB Writer: all items processed and saved.")
	}()

	wg.Wait()
	slog.Info("Process complete. All resources cleared cleanly.")
}
