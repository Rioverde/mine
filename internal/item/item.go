package item

import (
	"math/rand/v2"
	"time"

	"github.com/Rioverde/mine/internal/config"
	"github.com/Rioverde/mine/internal/ingot"
	"github.com/Rioverde/mine/internal/ore"
)

type Item struct {
	Ingot   *ingot.Ingot
	Quality int
}

func (it *Item) Name() string {
	switch it.Ingot.Ore.Material {
	case ore.Iron:
		return "Sword"
	case ore.Copper:
		return "Shield"
	case ore.Gold:
		return "Axe"
	}
	return "Unknown"
}

func NewSmithy(cfg config.ItemConfig) func(*ingot.Ingot) *Item {
	return func(in *ingot.Ingot) *Item {
		time.Sleep(500 * time.Millisecond)
		return &Item{
			Ingot:   in,
			Quality: min(in.Quality+rand.IntN(cfg.QualityBonus), cfg.MaxQuality),
		}
	}
}
