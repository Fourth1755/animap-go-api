package adapters

import (
	"net/http"

	"github.com/Fourth1755/animap-go-api/internal/core/dtos"
	"github.com/Fourth1755/animap-go-api/internal/core/services"
	"github.com/Fourth1755/animap-go-api/internal/errs"
	"github.com/gin-gonic/gin"
)

type HttpMyAnimeHandler struct {
	service services.MyAnimeService
}

func NewHttpMyAnimeHandler(service services.MyAnimeService) *HttpMyAnimeHandler {
	return &HttpMyAnimeHandler{service: service}
}

func (h *HttpMyAnimeHandler) AddAnimeToList(c *gin.Context) {
	MyAnimeRequest := new(dtos.AddAnimeToListRequest)
	if err := c.BindJSON(MyAnimeRequest); err != nil {
		handleError(c, errs.NewBadRequestError(err.Error()))
		return
	}

	if err := h.service.AddAnimeToList(MyAnimeRequest); err != nil {
		handleError(c, err)
		return
	}
	c.IndentedJSON(http.StatusOK, gin.H{"message": "Add anime to list success."})
}

func (h *HttpMyAnimeHandler) GetAnimeByUserId(c *gin.Context) {
	uuid := c.Param("uuid")

	animeList, err := h.service.GetAnimeByUserId(uuid)
	if err != nil {
		handleError(c, err)
		return
	}
	c.JSON(http.StatusOK, animeList)
}

func (h *HttpMyAnimeHandler) GetMyTopAnimeByUserId(c *gin.Context) {
	uuid := c.Param("uuid")

	animeList, err := h.service.GetMyTopAnime(uuid)
	if err != nil {
		handleError(c, err)
		return
	}
	c.JSON(http.StatusOK, animeList)
}

func (h *HttpMyAnimeHandler) UpdateMyTopAnime(c *gin.Context) {
	MyAnimeRequest := new(dtos.UpdateMyTopAnimeRequest)
	if err := c.BindJSON(MyAnimeRequest); err != nil {
		handleError(c, errs.NewBadRequestError(err.Error()))
		return
	}

	err := h.service.UpdateMyTopAnime(MyAnimeRequest)
	if err != nil {
		handleError(c, err)
		return
	}
	c.IndentedJSON(http.StatusOK, gin.H{"message": "Update my top anime to list success."})
}
