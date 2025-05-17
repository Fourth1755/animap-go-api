package repositories

import (
	"github.com/Fourth1755/animap-go-api/internal/core/entities"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type CategoryUniverse interface {
	GetAll() ([]entities.CategoryUniverse, error)
	GetById(id uuid.UUID) (*entities.CategoryUniverse, error)
}

type GormCategoryUniverse struct {
	db *gorm.DB
}

func NewGormCategoryUniverseRepository(db *gorm.DB) CategoryUniverse {
	return &GormCategoryUniverse{
		db: db,
	}
}

func (r *GormCategoryUniverse) GetAll() ([]entities.CategoryUniverse, error) {
	var categorise []entities.CategoryUniverse
	if result := r.db.Find(&categorise); result.Error != nil {
		return nil, result.Error
	}
	return categorise, nil
}

func (r *GormCategoryUniverse) GetById(id uuid.UUID) (*entities.CategoryUniverse, error) {
	category := new(entities.CategoryUniverse)
	if result := r.db.First(&category, id); result.Error != nil {
		return nil, result.Error
	}
	return category, nil
}
