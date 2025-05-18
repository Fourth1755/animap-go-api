package repositories

import (
	"github.com/Fourth1755/animap-go-api/internal/core/entities"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type StudioRepository interface {
	GetAll() ([]entities.Studio, error)
	GetByIds(ids []uuid.UUID) ([]entities.Studio, error)
}

type GormStudioRepository struct {
	db *gorm.DB
}

func NewGormStudioRepository(db *gorm.DB) StudioRepository {
	return &GormStudioRepository{db: db}
}

func (r GormStudioRepository) GetAll() ([]entities.Studio, error) {
	var studio []entities.Studio
	if result := r.db.Find(&studio); result.Error != nil {
		return nil, result.Error
	}
	return studio, nil
}

func (r GormStudioRepository) GetByIds(ids []uuid.UUID) ([]entities.Studio, error) {
	var studio []entities.Studio
	if result := r.db.Where("id in (?)", ids).
		Find(&studio); result.Error != nil {
		return nil, result.Error
	}
	return studio, nil
}
