package adapters

import (
	"net/http"

	"github.com/Fourth1755/animap-go-api/internal/core/dtos"
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
	//userAnime := new(entities.UserAnime)
	userAnimeRequest := new(dtos.AddAnimeToListRequest)
	if err := c.BindJSON(userAnimeRequest); err != nil {
		handleError(c, errs.NewBadRequestError(err.Error()))
		return
	}

	if err := h.service.AddAnimeToList(userAnimeRequest); err != nil {
		handleError(c, err)
		return
	}
	c.IndentedJSON(http.StatusOK, gin.H{"message": "Add anime to list success."})
}

func (h *HttpUserAnimeHandler) GetAnimeByUserId(c *gin.Context) {
	uuid := c.Param("uuid")

	animeList, err := h.service.GetAnimeByUserId(uuid)
	if err != nil {
		handleError(c, err)
		return
	}
	c.JSON(http.StatusOK, animeList)
}
