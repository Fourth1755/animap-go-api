package adapters

import (
	"net/http"
	"strconv"

	"github.com/Fourth1755/animap-go-api/internal/core/dtos"
	"github.com/Fourth1755/animap-go-api/internal/core/entities"
	"github.com/Fourth1755/animap-go-api/internal/core/services"
	"github.com/Fourth1755/animap-go-api/internal/errs"
	"github.com/gin-gonic/gin"
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
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		handleError(c, errs.NewBadRequestError(err.Error()))
		return
	}
	song, err := h.service.GetSongById(uint(id))
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
	songId, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		handleError(c, errs.NewBadRequestError(err.Error()))
		return
	}
	song := new(entities.Song)
	if err := c.BindJSON(&song); err != nil {
		handleError(c, err)
		return
	}
	song.ID = uint(songId)
	if err := h.service.UpdateSong(song); err != nil {
		handleError(c, err)
		return
	}
	c.IndentedJSON(http.StatusOK, gin.H{"message": "Update song success"})
}

func (h *HttpSongHandler) DeleteSong(c *gin.Context) {
	songId, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		handleError(c, errs.NewBadRequestError(err.Error()))
		return
	}
	if err := h.service.DeleteSong(uint(songId)); err != nil {
		handleError(c, err)
		return
	}
	c.IndentedJSON(http.StatusCreated, gin.H{"message": "Delete song success"})
}

func (h *HttpSongHandler) GetSongByAnimeId(c *gin.Context) {
	animeId, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		handleError(c, errs.NewBadRequestError(err.Error()))
		return
	}
	songs, err := h.service.GetSongByAnimeId(uint(animeId))
	if err != nil {
		handleError(c, err)
		return
	}
	c.JSON(http.StatusOK, songs)
}
