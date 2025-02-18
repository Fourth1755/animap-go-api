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

type SongChannel struct {
	ID      uint   `json:"id"`
	Channel int    `json:"channel"` // 1: youtube 2: spotify
	Type    int    `json:"type"`    // 1: tv_size 2: full 3: official 4 unofficial
	Link    string `json:"link"`
	IsMain  int    `json:"is_main"` // 1: main 2:secondary 3:not main is_main for show
}

type SongDetail struct {
	ID          uint          `json:"id"`
	Name        string        `json:"name"`
	Image       string        `json:"image"`
	Description string        `json:"description"`
	Year        string        `json:"year"`
	Type        int           `json:"type"` // 1: opening, 2: ending, 3:soundtrack
	Sequence    int           `json:"sequence"`
	AnimeID     uint          `json:"anime_id"`
	SongChannel []SongChannel `json:"song_channel"`
}
type SongDetailResponse struct {
	OpeningSong    []SongDetail `json:"opening_song"`
	EndingSong     []SongDetail `json:"ending_song"`
	SoundtrackSong []SongDetail `json:"soundtrack_song"`
}

type CreateSongRequest struct {
	Name        string        `json:"name"`
	Image       string        `json:"image"`
	Description string        `json:"description"`
	Year        string        `json:"year"`
	Type        int           `json:"type"` // 1: opening, 2: ending, 3:soundtrack
	Sequence    int           `json:"sequence"`
	AnimeID     uint          `json:"anime_id"`
	SongChannel []SongChannel `json:"song_channel"`
	ArtistList  []uint        `json:"artist_list"`
}

type CreateSongChannelRequest struct {
	Channel int    `json:"channel"`
	Type    int    `json:"type"`
	Link    string `json:"link"`
	SongID  uint   `json:"song_id"`
	IsMain  int    `json:"is_main"`
}
