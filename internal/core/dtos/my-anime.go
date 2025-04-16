package dtos

import (
	"time"

	"github.com/google/uuid"
)

type GetAnimeByUserIdResponse struct {
	AnimeID       uuid.UUID `json:"anime_id"`
	AnimeName     string    `json:"anime_name"`
	Score         float32   `json:"score"`
	Description   string    `json:"description"`
	Episodes      string    `json:"episodes"`
	Image         string    `json:"image"`
	Status        int       `json:"status"`
	WatchedAt     time.Time `json:"watched_at"`
	WatchedYearAt string    `json:"watched_year_at"`
	CreatedAt     time.Time `json:"create_at"`
}

type GetMyAnimeYearByUserIdResponseAnimeYearAnime struct {
	AnimeID       uuid.UUID `json:"anime_id"`
	AnimeName     string    `json:"anime_name"`
	Score         float32   `json:"score"`
	Description   string    `json:"description"`
	Episodes      string    `json:"episodes"`
	Image         string    `json:"image"`
	Status        int       `json:"status"`
	WatchedAt     time.Time `json:"watched_at"`
	WatchedYearAt string    `json:"watched_year_at"`
	CreatedAt     time.Time `json:"create_at"`
}

type GetMyAnimeYearByUserIdResponseAnimeYear struct {
	Year  string                                         `json:"year"`
	Anime []GetMyAnimeYearByUserIdResponseAnimeYearAnime `json:"anime"`
}

type GetMyAnimeYearByUserIdResponse struct {
	AnimeYear  []GetMyAnimeYearByUserIdResponseAnimeYear `json:"anime_year"`
	TotalYear  uint                                      `json:"total_year"`
	TotalAnime uint                                      `json:"total_anime"`
}

type AddAnimeToListRequest struct {
	UserUUID    uuid.UUID `json:"user_uuid"`
	AnimeID     uuid.UUID `json:"anime_id"`
	Score       float32   `json:"score"`
	Status      int       `json:"status"`
	WatchedYear string    `json:"watched_year"`
}

type GetMyTopAnimeResponse struct {
	AnimeID            uuid.UUID `json:"anime_id"`
	AnimeName          string    `json:"anime_name"`
	Score              float32   `json:"score"`
	Description        string    `json:"description"`
	Episodes           string    `json:"episodes"`
	Image              string    `json:"image"`
	Status             int       `json:"status"`
	WatchAt            time.Time `json:"watch_at"`
	CreatedAt          time.Time `json:"create_at"`
	SequenceMyTopAnime int       `json:"sequence_my_top_anime"`
}

type AnimeSequence struct {
	AnimeID  uuid.UUID `json:"anime_id"`
	Sequence int       `json:"sequence"`
}
type UpdateMyTopAnimeRequest struct {
	UserId        uuid.UUID       `json:"user_uuid"`
	AnimeSequence []AnimeSequence `json:"anime_sequence"`
}
