package repositories

import (
	"github.com/Fourth1755/animap-go-api/internal/core/entities"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type AnimeCharacterRepository interface {
	Save(category *entities.AnimeCharacter) error
	GetByAnimeId(animeId uuid.UUID) ([]entities.AnimeCharacter, error)
}

type GormAnimeCharacterRepository struct {
	dbPrimary *gorm.DB
	dbReplica *gorm.DB
}

func NewGormAnimeCharacterRepository(dbPrimary *gorm.DB, dbReplica *gorm.DB) AnimeCharacterRepository {
	return &GormAnimeCharacterRepository{dbPrimary: dbPrimary, dbReplica: dbReplica}
}

func (r *GormAnimeCharacterRepository) Save(character *entities.AnimeCharacter) error {
	if result := r.dbPrimary.Create(&character); result.Error != nil {
		return result.Error
	}
	return nil
}

func (r *GormAnimeCharacterRepository) GetByAnimeId(animeId uuid.UUID) ([]entities.AnimeCharacter, error) {
	var animeCharacter []entities.AnimeCharacter
	result := r.dbReplica.Preload("Character").Where("anime_id = ?", animeId).Find(&animeCharacter)
	if result.Error != nil {
		return nil, result.Error
	}
	return animeCharacter, nil
}
