package database

import (
	"fmt"

	"github.com/Fourth1755/animap-go-api/internal/core/config"
	"github.com/Fourth1755/animap-go-api/internal/core/entities"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func InitDatabase(cfgService config.ConfigService) (*gorm.DB, *gorm.DB) {
	dbConfig := cfgService.GetDatabase()
	dbConfigRepica := cfgService.GetDatabaseReplica()

	// Primary
	dsnPrimary := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable",
		dbConfig.Host,
		dbConfig.Port,
		dbConfig.UserName,
		dbConfig.Password,
		dbConfig.DatabaseName)
	dbPrimary, err := gorm.Open(postgres.Open(dsnPrimary), &gorm.Config{})
	if err != nil {
		panic("failed to connect primary database")
	}

	// Replica
	dsnReplica := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable",
		dbConfigRepica.Host,
		dbConfigRepica.Port,
		dbConfigRepica.UserName,
		dbConfigRepica.Password,
		dbConfigRepica.DatabaseName)
	dbReplica, err := gorm.Open(postgres.Open(dsnReplica), &gorm.Config{})
	if err != nil {
		panic("failed to connect replica database")
	}

	dbPrimary.AutoMigrate(
		&entities.Anime{},
		&entities.User{},
		&entities.UserAnime{},
		&entities.Category{},
		&entities.AnimeCategory{},
		&entities.Song{},
		&entities.SongChannel{},
		&entities.Artist{},
		&entities.SongArtist{},
		&entities.Studio{},
		&entities.AnimeStudio{},
		&entities.CategoryUniverse{},
		&entities.AnimeCategoryUniverse{},
		&entities.Episode{},
		&entities.Character{},
		&entities.AnimeCharacter{},
		&entities.EpisodeCharacter{},
		&entities.CommentAnime{},
		&entities.AnimePicture{},
		&entities.AnimeTrailer{},
	)

	return dbPrimary, dbReplica
}
