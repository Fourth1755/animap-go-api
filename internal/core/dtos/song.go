package dtos

type SongListResponse struct {
	ID          uint   `json:"id"`
	Name        string `json:"name"`
	Image       string `json:"image"`
	Description string `json:"description"`
	Year        string `json:"year"`
	Type        string `json:"type"` // 1: opening, 2: ending, 3:soundtrack
	Sequence    int    `json:"sequence"`
	AnimeID     uint   `json:"anime_id"`
	AnimeName   string `json:"anime_name"`
}
