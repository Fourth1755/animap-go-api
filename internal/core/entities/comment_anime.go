package entities

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type CommentAnime struct {
	ID           uuid.UUID `gorm:"primaryKey" json:"id"`
	Message      string
	AuthorID     uuid.UUID `gorm:"primaryKey" json:"author_id"`
	SubCommentID uuid.UUID `gorm:"primaryKey" json:"sub_comment_id"`
	IsEdit       bool
	Type         string
	AnimeID      uuid.UUID `gorm:"primaryKey" json:"anime_id"`
	CreatedAt    time.Time
	UpdatedAt    time.Time
	DeletedAt    gorm.DeletedAt `gorm:"index"`
}
