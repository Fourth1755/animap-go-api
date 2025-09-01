package entities

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Anime struct {
	ID                uuid.UUID  `gorm:"primarykey"`
	Name              string     `json:"name"`
	NameThai          string     `json:"name_thai"`
	NameEnglish       string     `json:"name_english"`
	Episodes          int        `json:"episodes"`
	Image             string     `json:"image"`
	Description       string     `json:"description"`
	Seasonal          string     `json:"seasonal"`
	Year              string     `json:"year"`
	Type              int        `json:"type"` //1: TV 2: movie
	MediaType         string     `json:"media_type"`
	Duration          string     `json:"duration"`
	Categories        []Category `gorm:"many2many:anime_categories;"`
	Songs             []Song
	Wallpaper         string             `json:"wallpaper"`
	Trailer           string             `json:"trailer"`
	TrailerEmbed      string             `json:"trailer_embed"`
	Studios           []Studio           `gorm:"many2many:anime_studios;"`
	CategoryUniverses []CategoryUniverse `gorm:"many2many:anime_category_universes;"`
	AiredAt           time.Time          `json:"aired_at"`
	CreatedAt         time.Time
	UpdatedAt         time.Time
	DeletedAt         gorm.DeletedAt `gorm:"index"`
	IsCreateEpisode   bool           `gorm:"default:false" json:"is_create_episode"`
	MyAnimeListID     uint64         `json:"my_anime_list_id"`
	IsSubAnime        bool           `json:"is_sub_anime"`
	Rating            string         `json:"rating"`
	IsShow            bool           `gorm:"default:true" json:"is_show"`
}
