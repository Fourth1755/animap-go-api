package repositories

import (
	"github.com/Fourth1755/animap-go-api/internal/core/entities"
	"gorm.io/gorm"
)

type SongChannelRepository interface {
	Save(songChannel *entities.SongChannel) error
	Update(songChannel *entities.SongChannel) error
}

type GormSongChannelRepository struct {
	dbPrimary *gorm.DB
	dbReplica *gorm.DB
}

func NewGormSongChannelRepository(dbPrimary *gorm.DB, dbReplica *gorm.DB) SongChannelRepository {
	return &GormSongChannelRepository{dbPrimary: dbPrimary, dbReplica: dbReplica}
}

func (r GormSongChannelRepository) Save(songChannel *entities.SongChannel) error {
	result := r.dbPrimary.Create(&songChannel)
	if result.Error != nil {
		return result.Error
	}
	return nil
}

func (r GormSongChannelRepository) Update(songChannel *entities.SongChannel) error {
	result := r.dbPrimary.Model(&songChannel).Updates(songChannel)
	if result.Error != nil {
		return result.Error
	}
	return nil
}
