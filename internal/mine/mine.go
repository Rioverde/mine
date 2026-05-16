package mine

import (
	"context"
	"log/slog"
	"time"

	"github.com/Rioverde/mine/internal/ore"
)

const (
	BufferSize = 20
	TickRate   = 100 * time.Millisecond
)

func Run(ctx context.Context) <-chan *ore.Ore {
	out := make(chan *ore.Ore, BufferSize)
	timer := time.NewTicker(TickRate)

	go func() {
		defer timer.Stop()
		defer close(out)
		for {
			select {
			case <-ctx.Done():
				slog.Info("Mine is closed due to the timer")
				return
			case <-timer.C:
				o := ore.Random()
				select {
				case out <- o:
					slog.Info("Ore found", "material", o.Name(), "capacity", o.Capacity)
				case <-ctx.Done():
					slog.Info("Mine is closed due to the timer")
					return
				}
			}
		}
	}()

	return out
}
