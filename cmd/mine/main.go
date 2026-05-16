package main

import (
	"context"
	"log"
	"log/slog"
	"os/signal"
	"syscall"

	"github.com/Rioverde/mine/internal/config"
	"github.com/Rioverde/mine/internal/factory"
	"github.com/Rioverde/mine/internal/ingot"
	"github.com/Rioverde/mine/internal/item"
	"github.com/Rioverde/mine/internal/mine"
	"github.com/Rioverde/mine/internal/storage"
)

func main() {
	cfg, err := config.ReadConfig("config.yaml")
	if err != nil {
		log.Fatal(err)
	}

	ctx, cancel := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer cancel()

	pool, err := storage.Connect(ctx)
	if err != nil {
		log.Fatal(err)
	}
	defer pool.Close()

	ores := mine.Run(ctx, cfg)
	ingots := factory.New(factory.WithName(ctx, "Smelter"), cfg, ores, ingot.NewSmelter(cfg.Ingot))
	items := factory.New(factory.WithName(ctx, "Smithy"), cfg, ingots, item.NewSmithy(cfg.Item))

	for v := range items {
		id, err := storage.SaveItem(ctx, pool, "Smithy", v)
		if err != nil {
			slog.Error("save item", "err", err)
			continue
		}
		slog.Info("Weapon delivered",
			"id", id,
			"ore", v.Ingot.Ore.Name(),
			"type", v.Name(),
			"quality", v.Quality,
		)
	}

	slog.Info("Process complete")
}
