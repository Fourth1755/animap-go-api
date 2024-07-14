package entities

import "gorm.io/gorm"

type AnimeCategory struct {
	gorm.Model
	CategoryID       uint `json:"category_id"`
	Category         Category
	AnimeID          uint `json:"anime_id"`
	Anime            Anime
	Description      string `json:"description"`
	Sequence         int    `json:"sequence"`
	SequenceTimeLine int    `json:"sequence_time_line"`
}
