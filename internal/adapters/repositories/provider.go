package repositories

import (
	"github.com/Fourth1755/animap-go-api/internal/core/entities"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type ProviderRepository interface {
	Save(provider *entities.Provider) (*entities.Provider, error)
	GetAll() ([]entities.Provider, error)
	GetById(id uuid.UUID) (*entities.Provider, error)
}

type GormProviderRepository struct {
	dbPrimary *gorm.DB
	dbReplica *gorm.DB
}

func NewGormProviderRepository(dbPrimary *gorm.DB, dbReplica *gorm.DB) ProviderRepository {
	return &GormProviderRepository{dbPrimary: dbPrimary, dbReplica: dbReplica}
}

func (r *GormProviderRepository) Save(provider *entities.Provider) (*entities.Provider, error) {
	if result := r.dbPrimary.Create(provider); result.Error != nil {
		return nil, result.Error
	}
	return provider, nil
}

func (r *GormProviderRepository) GetAll() ([]entities.Provider, error) {
	var providers []entities.Provider
	if result := r.dbReplica.Find(&providers); result.Error != nil {
		return nil, result.Error
	}
	return providers, nil
}

func (r *GormProviderRepository) GetById(id uuid.UUID) (*entities.Provider, error) {
	provider := new(entities.Provider)
	if result := r.dbReplica.First(provider, id); result.Error != nil {
		return nil, result.Error
	}
	return provider, nil
}
