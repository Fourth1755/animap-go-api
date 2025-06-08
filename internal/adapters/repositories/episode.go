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
	db *gorm.DB
}

func NewGormEpisodeRepository(db *gorm.DB) EpisodeRepository {
	return &GormEpisodeRepository{db: db}
}

func (r *GormEpisodeRepository) BulkSave(episodes []entities.Episode) error {
	if result := r.db.Create(&episodes); result.Error != nil {
		return result.Error
	}
	return nil
}

func (r *GormEpisodeRepository) GetByAnimeId(anime_id uuid.UUID) ([]entities.Episode, error) {
	var animeEpisode []entities.Episode
	result := r.db.Where("anime_id = ?", anime_id).Find(&animeEpisode)
	if result.Error != nil {
		return nil, result.Error
	}
	return animeEpisode, nil
}

func (r *GormEpisodeRepository) Update(animeEpisode *entities.Episode) error {
	result := r.db.Model(&animeEpisode).Updates(animeEpisode)
	if result.Error != nil {
		return result.Error
	}
	return nil
}

func (r *GormEpisodeRepository) GetById(id uuid.UUID) (*entities.Episode, error) {
	var episode entities.Episode
	if result := r.db.
		First(&episode, id); result.Error != nil {
		return nil, result.Error
	}
	return &episode, nil
}
