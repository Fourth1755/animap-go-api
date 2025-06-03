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

func (h *HttpEpisodeHandler) GetByAnimeId(c *gin.Context) {
	animeId := c.Param("anime_id")
	anime, err := h.service.GetByAnimeId(uuid.MustParse(animeId))
	if err != nil {
		handleError(c, err)
		return
	}
	c.JSON(http.StatusOK, anime)
}
