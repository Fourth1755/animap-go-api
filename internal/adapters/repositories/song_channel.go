package repositories

import (
	"github.com/Fourth1755/animap-go-api/internal/core/entities"
	"gorm.io/gorm"
)

type SongChannelRepository interface {
	Save(songChannel *entities.SongChannel) error
}

type GormSongChannelRepository struct {
	db *gorm.DB
}

func NewGormSongChannelRepository(db *gorm.DB) SongChannelRepository {
	return &GormSongChannelRepository{db: db}
}

func (r GormSongChannelRepository) Save(songChannel *entities.SongChannel) error {
	result := r.db.Create(&songChannel)
	if result.Error != nil {
		return result.Error
	}
	return nil
}
