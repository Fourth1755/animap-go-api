package dtos

import "time"

type GetAnimeByUserIdResponse struct {
	AnimeID       uint      `json:"anime_id"`
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

type GetMyAnimeYearByUserIdResponse_AnimeYear_Anime struct {
	AnimeID       uint      `json:"anime_id"`
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

type GetMyAnimeYearByUserIdResponse_AnimeYear struct {
	Year  string                                           `json:"year"`
	Anime []GetMyAnimeYearByUserIdResponse_AnimeYear_Anime `json:"anime"`
}

type GetMyAnimeYearByUserIdResponse struct {
	AnimeYear  []GetMyAnimeYearByUserIdResponse_AnimeYear `json:"anime_year"`
	TotalYear  uint                                       `json:"total_year"`
	TotalAnime uint                                       `json:"total_anime"`
}

type AddAnimeToListRequest struct {
	UserUUID    string  `json:"user_uuid"`
	AnimeID     uint    `json:"anime_id"`
	Score       float32 `json:"score"`
	Status      int     `json:"status"`
	WatchedYear string  `json:"watched_year"`
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
