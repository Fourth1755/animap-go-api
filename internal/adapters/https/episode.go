package adapters

import (
	"net/http"

	"github.com/Fourth1755/animap-go-api/internal/core/dtos"
	"github.com/Fourth1755/animap-go-api/internal/core/services"
	"github.com/Fourth1755/animap-go-api/internal/errs"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type HttpEpisodeHandler struct {
	service services.EpisodeService
}

func NewHttpEpisodeHandler(service services.EpisodeService) *HttpEpisodeHandler {
	return &HttpEpisodeHandler{service: service}
}

func (h *HttpEpisodeHandler) CreateEpisode(c *gin.Context) {
	var request dtos.CreateEpisodeRequest
	if err := c.BindJSON(&request); err != nil {
		handleError(c, errs.NewBadRequestError(err.Error()))
		return
	}
	err := h.service.CreateEpisode(request.AnimeId)
	if err != nil {
		handleError(c, err)
		return
	}
	c.IndentedJSON(http.StatusCreated, gin.H{"message": "Create Episode success"})
}

func (h *HttpEpisodeHandler) GetEpisodesByAnimeId(c *gin.Context) {
	filter := c.Query("filter")
	animeId := c.Param("anime_id")
	anime, err := h.service.GetEpisodesByAnimeId(uuid.MustParse(animeId), filter)
	if err != nil {
		handleError(c, err)
		return
	}
	c.JSON(http.StatusOK, anime)
}

func (h *HttpEpisodeHandler) UpdateEpisode(c *gin.Context) {
	request := new(dtos.UpdateEpisodeRequest)
	if err := c.BindJSON(request); err != nil {
		handleError(c, errs.NewBadRequestError(err.Error()))
		return
	}
	err := h.service.UpdateEpisode(*request)
	if err != nil {
		handleError(c, errs.NewBadRequestError(err.Error()))
		return
	}
	c.IndentedJSON(http.StatusOK, gin.H{"message": "Update Episode success"})
}

func (h *HttpEpisodeHandler) AddCharactersToEpisode(c *gin.Context) {
	request := new(dtos.AddCharacterToEpisodeRequest)
	if err := c.BindJSON(request); err != nil {
		handleError(c, errs.NewBadRequestError(err.Error()))
		return
	}
	err := h.service.AddCharactersToEpisode(*request)
	if err != nil {
		handleError(c, errs.NewBadRequestError(err.Error()))
		return
	}
	c.IndentedJSON(http.StatusOK, gin.H{"message": "Add Characters to Episode success"})
}
