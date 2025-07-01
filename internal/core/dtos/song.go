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
	Link    string    `json:"link"`
	Channel string    `json:"channel"` // YOUTUBE,SPOTIFY
	Type    string    `json:"type"`    // TV_SIZE, FULL_SIZE_OFFICIAL, FULL_SIZE_UNOFFICIAL, FIRST_TAKE
	IsMain  bool      `json:"is_main"` // true: main false:not main is_main for show
}

type GetSongByAnimeIdResponseSong struct {
	ID          uuid.UUID                             `json:"id"`
	Name        string                                `json:"name"`
	Image       string                                `json:"image"`
	Description string                                `json:"description"`
	Year        string                                `json:"year"`
	Type        string                                `json:"type"` //OPENING, ENDING, SOUNDSTRACK
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
	Type        string                                `json:"type"` //OPENING, ENDING, SOUNDSTRACK
	Sequence    int                                   `json:"sequence"`
	AnimeID     uuid.UUID                             `json:"anime_id"`
	SongChannel []GetSongByAnimeIdResponseSongChannel `json:"song_channel"`
	ArtistList  []uuid.UUID                           `json:"artist_list"`
}

type CreateSongChannelRequest struct {
	Channel string    `json:"channel"`
	Type    string    `json:"type"`
	Link    string    `json:"link"`
	SongID  uuid.UUID `json:"song_id"`
	IsMain  bool      `json:"is_main"`
}

type GetSongsByArtistResponseSongChannel struct {
	ID      uuid.UUID `json:"id"`
	Link    string    `json:"link"`
	Channel string    `json:"channel"` // YOUTUBE,SPOTIFY
	Type    string    `json:"type"`    // TV_SIZE, FULL_SIZE_OFFICIAL, FULL_SIZE_UNOFFICIAL, FIRST_TAKE
	IsMain  bool      `json:"is_main"` // true: main false:not main is_main for show
}
type GetSongsByArtistResponseSong struct {
	ID             uuid.UUID                             `json:"id"`
	Name           string                                `json:"name"`
	Image          string                                `json:"image"`
	Description    string                                `json:"description"`
	Year           string                                `json:"year"`
	Type           string                                `json:"type"` //OPENING, ENDING, SOUNDSTRACK
	AnimeID        uuid.UUID                             `json:"anime_id"`
	AnimeName      string                                `json:"anime_name"`
	AnimeWallpaper string                                `json:"anime_wallpaper"`
	SongChannel    []GetSongsByArtistResponseSongChannel `json:"song_channel"`
}

type GetSongsByArtistResponseArtist struct {
	ID    uuid.UUID `json:"id"`
	Name  string    `json:"name"`
	Image string    `json:"image"`
}

type GetSongsByArtistResponse struct {
	Aritst GetSongsByArtistResponseArtist `json:"aritst"`
	Songs  []GetSongsByArtistResponseSong `json:"songs"`
}

type AddSongChannelsToSongRequest struct {
	SongID  uuid.UUID `json:"song_id"`
	Link    string    `json:"link"`
	Channel string    `json:"channel"` // YOUTUBE,SPOTIFY
	Type    string    `json:"type"`    // TV_SIZE, FULL_SIZE_OFFICIAL, FULL_SIZE_UNOFFICIAL, FIRST_TAKE
	IsMain  bool      `json:"is_main"` // true: main false:not main is_main for show
}
