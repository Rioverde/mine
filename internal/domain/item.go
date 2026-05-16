package domain

import (
	"math/rand/v2"
	"time"

	"github.com/Rioverde/mine/internal/config"
)

type Item struct {
	Ingot   *Ingot
	Quality int
}

func (it *Item) Name() string {
	switch it.Ingot.Ore.Material {
	case Iron:
		return "Sword"
	case Copper:
		return "Shield"
	case Gold:
		return "Axe"
	}
	return "Unknown"
}

func NewSmithy(cfg config.ItemConfig) func(*Ingot) *Item {
	return func(in *Ingot) *Item {
		time.Sleep(500 * time.Millisecond)
		return &Item{
			Ingot:   in,
			Quality: min(in.Quality+rand.IntN(cfg.QualityBonus), cfg.MaxQuality),
		}
	}
}
