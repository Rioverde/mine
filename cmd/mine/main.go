package main

import (
	"context"
	"log"
	"log/slog"
	"os"
	"time"

	"github.com/Rioverde/mine/internal/factory"
	"github.com/Rioverde/mine/internal/ingot"
	"github.com/Rioverde/mine/internal/item"
	"github.com/Rioverde/mine/internal/mine"
)

const (
	runDuration     = 5 * time.Second
	storageCapacity = 100
	logFile         = "./app.log"
)

func main() {
	f, err := os.OpenFile(logFile, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0o666)
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()

	slog.SetDefault(slog.New(slog.NewTextHandler(f, nil)))

	ctx, cancel := context.WithTimeout(context.Background(), runDuration)
	defer cancel()

	ores := mine.Run(ctx)
	ingots := factory.New(factory.WithName(ctx, "Smelter"), ores, ingot.FromOre)
	items := factory.New(factory.WithName(ctx, "Smithy"), ingots, item.FromIngot)

	storage := make([]*item.Item, 0, storageCapacity)
	for v := range items {
		slog.Info("Weapon delivered",
			"ore", v.Ingot.Ore.Name(),
			"type", v.Name(),
			"quality", v.Quality,
		)
		storage = append(storage, v)
	}

	slog.Info("Process complete", "items", len(storage))
}
