package repositories

import (
	"github.com/Fourth1755/animap-go-api/internal/core/entities"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type CategoryUniverseRepository interface {
	Save(category *entities.CategoryUniverse) error
	GetAll() ([]entities.CategoryUniverse, error)
	GetById(id uuid.UUID) (*entities.CategoryUniverse, error)
}

type GormCategoryUniverseRepository struct {
	dbPrimary *gorm.DB
	dbReplica *gorm.DB
}

func NewGormCategoryUniverseRepository(dbPrimary *gorm.DB, dbReplica *gorm.DB) CategoryUniverseRepository {
	return &GormCategoryUniverseRepository{
		dbPrimary: dbPrimary,
		dbReplica: dbReplica,
	}
}

func (r *GormCategoryUniverseRepository) Save(category *entities.CategoryUniverse) error {
	if result := r.dbPrimary.Create(&category); result.Error != nil {
		return result.Error
	}
	return nil
}

func (r *GormCategoryUniverseRepository) GetAll() ([]entities.CategoryUniverse, error) {
	var categorise []entities.CategoryUniverse
	if result := r.dbReplica.Find(&categorise); result.Error != nil {
		return nil, result.Error
	}
	return categorise, nil
}

func (r *GormCategoryUniverseRepository) GetById(id uuid.UUID) (*entities.CategoryUniverse, error) {
	category := new(entities.CategoryUniverse)
	if result := r.dbReplica.First(&category, id); result.Error != nil {
		return nil, result.Error
	}
	return category, nil
}
