package dtos

type AnimeDTO struct {
	ID       uint   `json:"id"`
	Name     string `json:"name"`
	Episodes int    `json:"episodes"`
	Seasonal string `json:"seasonal"`
	Year     string `json:"year"`
}
type AnimeDetailCategories struct {
	ID   uint   `json:"id"`
	Name string `json:"name"`
}

type AnimeDetailSongs struct {
	ID       uint   `json:"id"`
	Name     string `json:"name"`
	Type     uint   `json:"type"`
	Sequence uint   `json:"sequence"`
	Link     string `json:"link"`
}
type AnimeDetailResponse struct {
	ID             uint                    `json:"id"`
	Name           string                  `json:"name"`
	NameEnglish    string                  `json:"name_english"`
	Episodes       int                     `json:"episodes"`
	Seasonal       string                  `json:"seasonal"`
	Year           string                  `json:"year"`
	Image          string                  `json:"image"`
	Description    string                  `json:"description"`
	Duration       string                  `json:"duration"`
	Type           int                     `json:"type"`
	Categories     []AnimeDetailCategories `json:"categories"`
	OpeningSong    []AnimeDetailSongs      `json:"opening"`
	EndingSong     []AnimeDetailSongs      `json:"ending"`
	SoundtrackSong []AnimeDetailSongs      `json:"soundtrack"`
}

type AnimeQueryDTO struct {
	Seasonal string
	Year     string
}
