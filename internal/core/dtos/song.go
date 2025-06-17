package dtos

import "github.com/google/uuid"

type SongListResponse struct {
	ID          uuid.UUID `json:"id"`
	Name        string    `json:"name"`
	Image       string    `json:"image"`
	Description string    `json:"description"`
	Year        string    `json:"year"`
	Type        string    `json:"type"` // 1: opening, 2: ending, 3:soundtrack
	Sequence    int       `json:"sequence"`
	AnimeID     uuid.UUID `json:"anime_id"`
	AnimeName   string    `json:"anime_name"`
}

type GetSongByAnimeIdResponseSongArtist struct {
	ID    uuid.UUID `json:"id"`
	Name  string    `json:"name"`
	Image string    `json:"image"`
}

type GetSongByAnimeIdResponseSongChannel struct {
	ID      uuid.UUID `json:"id"`
	Channel int       `json:"channel"` // 1: youtube 2: spotify
	Type    int       `json:"type"`    // 1: tv_size 2: full 3: official 4 unofficial
	Link    string    `json:"link"`
	IsMain  int       `json:"is_main"` // 1: main 2:secondary 3:not main is_main for show
}

type GetSongByAnimeIdResponseSong struct {
	ID          uuid.UUID                             `json:"id"`
	Name        string                                `json:"name"`
	Image       string                                `json:"image"`
	Description string                                `json:"description"`
	Year        string                                `json:"year"`
	Type        int                                   `json:"type"` // 1: opening, 2: ending, 3:soundtrack
	Sequence    int                                   `json:"sequence"`
	AnimeID     uuid.UUID                             `json:"anime_id"`
	SongArtist  []GetSongByAnimeIdResponseSongArtist  `json:"song_artist"`
	SongChannel []GetSongByAnimeIdResponseSongChannel `json:"song_channel"`
}
type GetSongByAnimeIdResponse struct {
	OpeningSong    []GetSongByAnimeIdResponseSong `json:"opening_song"`
	EndingSong     []GetSongByAnimeIdResponseSong `json:"ending_song"`
	SoundtrackSong []GetSongByAnimeIdResponseSong `json:"soundtrack_song"`
}

type CreateSongRequest struct {
	Name        string                                `json:"name"`
	Image       string                                `json:"image"`
	Description string                                `json:"description"`
	Year        string                                `json:"year"`
	Type        int                                   `json:"type"` // 1: opening, 2: ending, 3:soundtrack
	Sequence    int                                   `json:"sequence"`
	AnimeID     uuid.UUID                             `json:"anime_id"`
	SongChannel []GetSongByAnimeIdResponseSongChannel `json:"song_channel"`
	ArtistList  []uuid.UUID                           `json:"artist_list"`
}

type CreateSongChannelRequest struct {
	Channel int       `json:"channel"`
	Type    int       `json:"type"`
	Link    string    `json:"link"`
	SongID  uuid.UUID `json:"song_id"`
	IsMain  int       `json:"is_main"`
}
