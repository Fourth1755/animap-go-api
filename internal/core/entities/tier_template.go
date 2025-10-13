package entities

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type TierTemplate struct {
	ID          uuid.UUID `gorm:"primarykey"`
	Name        string    `json:"name"`
	Type        string
	PlayedCount uint
	TierList    map[string]interface{} `gorm:"serializer:json"`
	TotalItem   uint
	CreatedAt   time.Time
	UpdatedAt   time.Time
	DeletedAt   gorm.DeletedAt `gorm:"index"`
}
