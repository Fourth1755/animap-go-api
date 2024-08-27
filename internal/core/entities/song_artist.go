package entities

import "gorm.io/gorm"

type SongArtist struct {
	gorm.Model
	SongId   uint `json:"song_id"`
	ArtistId uint `json:"artist_id"`
}
