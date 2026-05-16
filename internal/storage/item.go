package storage

import (
	"fmt"
	"time"
)

type ItemDTO struct {
	id        string
	name      string
	factory   string
	material  string
	quality   int
	createdAt time.Time
}

func (dto *ItemDTO) ToJson() string {
	return fmt.Sprintf(`{"id":"%s","name":"%s","material":"%s","quality":%d,"factory":"%s","created_at":"%s"}`,
		dto.id,
		dto.name,
		dto.material,
		dto.quality,
		dto.factory,
		dto.createdAt.Format(time.RFC3339),
	)
}
func (dto *ItemDTO) String() string {
	return fmt.Sprintf("ID: %s, Name: %s, Material: %s, Quality: %d, Factory: %s, CreatedAt: %s",
		dto.id,
		dto.name,
		dto.material,
		dto.quality,
		dto.factory,
		dto.createdAt.Format("2006-01-02 15:04:05"),
	)
}
