package events

import (
	"strconv"
	"strings"

	"github.com/Rioverde/mine/internal/domain"
	"github.com/go-viper/mapstructure/v2"
)

type ItemPayload struct {
	ID       string `mapstructure:"id"`
	Name     string `mapstructure:"name"`
	Material string `mapstructure:"material"`
	Quality  int    `mapstructure:"quality"`
	Factory  string `mapstructure:"factory"`
}

func NewItemPayload(id, factory string, it *domain.Item) ItemPayload {
	return ItemPayload{
		ID:       id,
		Name:     it.Name(),
		Material: strings.ToLower(it.Ingot.Ore.Name()),
		Quality:  it.Quality,
		Factory:  factory,
	}
}

func (p ItemPayload) ToMap() map[string]any {
	return map[string]any{
		"id":       p.ID,
		"name":     p.Name,
		"material": p.Material,
		"quality":  strconv.Itoa(p.Quality),
		"factory":  p.Factory,
	}
}

func ParseItemPayload(ev Event) (ItemPayload, error) {
	var p ItemPayload
	err := mapstructure.WeakDecode(ev.Values, &p)
	return p, err
}
