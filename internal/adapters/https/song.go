package adapters

import (
	"net/http"

	"github.com/Fourth1755/animap-go-api/internal/core/dtos"
	"github.com/Fourth1755/animap-go-api/internal/core/entities"
	"github.com/Fourth1755/animap-go-api/internal/core/services"
	"github.com/Fourth1755/animap-go-api/internal/errs"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type HttpSongHandler struct {
	service services.SongService
}

func NewHttpSongHandler(service services.SongService) *HttpSongHandler {
	return &HttpSongHandler{service: service}
}

func (h *HttpSongHandler) CreateSong(c *gin.Context) {
	var song dtos.CreateSongRequest
	if err := c.BindJSON(&song); err != nil {
		handleError(c, errs.NewBadRequestError(err.Error()))
		return
	}
	if err := h.service.CreateSong(&song); err != nil {
		handleError(c, err)
		return
	}
	c.IndentedJSON(http.StatusCreated, gin.H{"message": "Create song success"})
}

func (h *HttpSongHandler) GetSongById(c *gin.Context) {
	songId := c.Param("id")
	song, err := h.service.GetSongById(uuid.MustParse(songId))
	if err != nil {
		handleError(c, err)
		return
	}
	c.JSON(http.StatusOK, song)
}

func (h *HttpSongHandler) GetSongAll(c *gin.Context) {
	song, err := h.service.GetAllSongs()
	if err != nil {
		handleError(c, err)
		return
	}
	c.JSON(http.StatusOK, song)
}

func (h *HttpSongHandler) UpdateSong(c *gin.Context) {
	songId := c.Param("id")
	song := new(entities.Song)
	if err := c.BindJSON(&song); err != nil {
		handleError(c, err)
		return
	}
	song.ID = uuid.MustParse(songId)
	if err := h.service.UpdateSong(song); err != nil {
		handleError(c, err)
		return
	}
	c.IndentedJSON(http.StatusOK, gin.H{"message": "Update song success"})
}

func (h *HttpSongHandler) DeleteSong(c *gin.Context) {
	songId := c.Param("id")
	if err := h.service.DeleteSong(uuid.MustParse(songId)); err != nil {
		handleError(c, err)
		return
	}
	c.IndentedJSON(http.StatusCreated, gin.H{"message": "Delete song success"})
}

func (h *HttpSongHandler) GetSongByAnimeId(c *gin.Context) {
	animeId := c.Param("id")
	songs, err := h.service.GetSongByAnimeId(uuid.MustParse(animeId))
	if err != nil {
		handleError(c, err)
		return
	}
	c.JSON(http.StatusOK, songs)
}

func (h *HttpSongHandler) CreateSongChannel(c *gin.Context) {
	var songChannel dtos.CreateSongChannelRequest
	if err := c.BindJSON(&songChannel); err != nil {
		handleError(c, errs.NewBadRequestError(err.Error()))
		return
	}
	if err := h.service.CreateSongChannel(&songChannel); err != nil {
		handleError(c, err)
		return
	}
	c.IndentedJSON(http.StatusCreated, gin.H{"message": "Create song chennel success"})
}

func (h *HttpSongHandler) GetSongsByArtistId(c *gin.Context) {
	artistId := c.Param("id")
	rersponse, err := h.service.GetSongsByArtistId(uuid.MustParse(artistId))
	if err != nil {
		handleError(c, err)
		return
	}
	c.JSON(http.StatusOK, rersponse)
}
