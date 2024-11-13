package dtos

import "time"

type UserAnimeListDTO struct {
	AnimeID     uint      `json:"anime_id"`
	AnimeName   string    `json:"anime_name"`
	Score       float32   `json:"score"`
	Description string    `json:"description"`
	Episodes    string    `json:"episodes"`
	Image       string    `json:"image"`
	Status      int       `json:"status"`
	WatchAt     time.Time `json:"watch_at"`
	CreatedAt   time.Time `json:"create_at"`
}

type AddAnimeToListRequest struct {
	UserUUID string  `json:"user_uuid"`
	AnimeID  uint    `json:"anime_id"`
	Score    float32 `json:"score"`
	Status   int     `json:"status"`
}

type GetMyTopAnimeResponse struct {
	AnimeID            uint      `json:"anime_id"`
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
	AnimeID  uint `json:"anime_id"`
	Sequence int  `json:"sequence"`
}
type UpdateMyTopAnimeRequest struct {
	UserUUID      string          `json:"user_uuid"`
	AnimeSequence []AnimeSequence `json:"anime_sequence"`
}
