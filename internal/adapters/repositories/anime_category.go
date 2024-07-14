package repositories

import (
	"github.com/Fourth1755/animap-go-api/internal/core/entities"
	"github.com/Fourth1755/animap-go-api/internal/core/ports"
	"gorm.io/gorm"
)

type GormAnimeCategoryRepository struct {
	db *gorm.DB
}

func NewGormAnimeCategoryRepository(db *gorm.DB) ports.AnimeCategoryRepository {
	return &GormAnimeCategoryRepository{db: db}
}

func (r GormAnimeCategoryRepository) Save(animeCategory *entities.AnimeCategory) error {
	if result := r.db.Create(&animeCategory); result.Error != nil {
		return result.Error
	}
	return nil
}
