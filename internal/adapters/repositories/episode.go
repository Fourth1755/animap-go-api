package repositories

import (
	"github.com/Fourth1755/animap-go-api/internal/core/entities"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type EpisodeRepository interface {
	BulkSave(episodes []entities.Episode) error
	GetByAnimeId(anime_id uuid.UUID) ([]entities.Episode, error)
	Update(animeEpisode *entities.Episode) error
	GetById(id uuid.UUID) (*entities.Episode, error)
}

type GormEpisodeRepository struct {
	dbPrimary *gorm.DB
	dbReplica *gorm.DB
}

func NewGormEpisodeRepository(dbPrimary *gorm.DB, dbReplica *gorm.DB) EpisodeRepository {
	return &GormEpisodeRepository{dbPrimary: dbPrimary, dbReplica: dbReplica}
}

func (r *GormEpisodeRepository) BulkSave(episodes []entities.Episode) error {
	if result := r.dbPrimary.Create(&episodes); result.Error != nil {
		return result.Error
	}
	return nil
}

func (r *GormEpisodeRepository) GetByAnimeId(anime_id uuid.UUID) ([]entities.Episode, error) {
	var episode []entities.Episode
	result := r.dbReplica.
		Order("number asc").
		Where("anime_id = ?", anime_id).Find(&episode)
	if result.Error != nil {
		return nil, result.Error
	}
	return episode, nil
}

func (r *GormEpisodeRepository) Update(episode *entities.Episode) error {
	result := r.dbPrimary.Model(&episode).Updates(episode)
	if result.Error != nil {
		return result.Error
	}
	return nil
}

func (r *GormEpisodeRepository) GetById(id uuid.UUID) (*entities.Episode, error) {
	var episode entities.Episode
	if result := r.dbReplica.
		First(&episode, id); result.Error != nil {
		return nil, result.Error
	}
	return &episode, nil
}
