package repositories

import (
	"github.com/Fourth1755/animap-go-api/internal/core/entities"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type AnimeProviderRepository interface {
	Save(animeProviders []entities.AnimeProvider) error
	GetByAnimeIdAndProviderIds(animeID uuid.UUID, providerIDs []uuid.UUID) ([]entities.AnimeProvider, error)
}

type GormAnimeProviderRepository struct {
	dbPrimary *gorm.DB
	dbReplica *gorm.DB
}

func NewGormAnimeProviderRepository(dbPrimary *gorm.DB, dbReplica *gorm.DB) AnimeProviderRepository {
	return &GormAnimeProviderRepository{dbPrimary: dbPrimary, dbReplica: dbReplica}
}

func (r *GormAnimeProviderRepository) Save(animeProviders []entities.AnimeProvider) error {
	if result := r.dbPrimary.Create(&animeProviders); result.Error != nil {
		return result.Error
	}
	return nil
}

func (r *GormAnimeProviderRepository) GetByAnimeIdAndProviderIds(animeID uuid.UUID, providerIDs []uuid.UUID) ([]entities.AnimeProvider, error) {
	var animeProviders []entities.AnimeProvider
	result := r.dbReplica.
		Where("anime_id = ? AND provider_id IN (?)", animeID, providerIDs).
		Find(&animeProviders)
	if result.Error != nil {
		return nil, result.Error
	}
	return animeProviders, nil
}
