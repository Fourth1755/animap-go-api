package adapters

import (
	"net/http"
	"strconv"

	"github.com/Fourth1755/animap-go-api/internal/core/entities"
	"github.com/Fourth1755/animap-go-api/internal/core/services"
	"github.com/Fourth1755/animap-go-api/internal/errs"
	"github.com/gin-gonic/gin"
)

type HttpUserAnimeHandler struct {
	service services.UserAnimeService
}

func NewHttpUserAnimeHandler(service services.UserAnimeService) *HttpUserAnimeHandler {
	return &HttpUserAnimeHandler{service: service}
}

func (h *HttpUserAnimeHandler) AddAnimeToList(c *gin.Context) {
	userAnime := new(entities.UserAnime)
	if err := c.BindJSON(userAnime); err != nil {
		handleError(c, errs.NewBadRequestError(err.Error()))
		return
	}

	if err := h.service.AddAnimeToList(userAnime); err != nil {
		handleError(c, err)
		return
	}
	c.IndentedJSON(http.StatusOK, gin.H{"message": "Add anime to list success."})
}

func (h *HttpUserAnimeHandler) GetAnimeByUserId(c *gin.Context) {
	userId, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		handleError(c, errs.NewBadRequestError(err.Error()))
		return
	}

	animeList, err := h.service.GetAnimeByUserId(uint(userId))
	if err != nil {
		handleError(c, err)
		return
	}
	c.JSON(http.StatusOK, animeList)
}
