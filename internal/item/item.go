package item

import (
	"math/rand/v2"
	"time"

	"github.com/Rioverde/mine/internal/config"
	"github.com/Rioverde/mine/internal/ingot"
)

const (
	Axe = iota
	Sword
	Shield
)

type Item struct {
	ItemType int
	Ingot    *ingot.Ingot
	Quality  int
}

func (it *Item) Name() string {
	switch it.ItemType {
	case Sword:
		return "Sword"
	case Shield:
		return "Shield"
	case Axe:
		return "Axe"
	}
	return "Unknown"
}

func NewSmithy(cfg config.ItemConfig) func(*ingot.Ingot) *Item {
	return func(in *ingot.Ingot) *Item {
		time.Sleep(500 * time.Millisecond)
		return &Item{
			ItemType: rand.IntN(cfg.NumberOfItems),
			Ingot:    in,
			Quality:  min(in.Quality+rand.IntN(cfg.QualityBonus), cfg.MaxQuality),
		}
	}
}
