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

type HttpAnimeHandler struct {
	service services.AnimeService
}

func NewHttpAnimeHandler(service services.AnimeService) *HttpAnimeHandler {
	return &HttpAnimeHandler{service: service}
}

func (h *HttpAnimeHandler) CreateAnime(c *gin.Context) {
	var anime dtos.CreateAnimeRequest
	if err := c.BindJSON(&anime); err != nil {
		handleError(c, errs.NewBadRequestError(err.Error()))
		return
	}
	if err := h.service.CreateAnime(anime); err != nil {
		handleError(c, err)
		return
	}
	c.IndentedJSON(http.StatusCreated, gin.H{"message": "Create anime success"})
}

func (h *HttpAnimeHandler) GetAnimeById(c *gin.Context) {
	animeId, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		handleError(c, errs.NewBadRequestError(err.Error()))
		return
	}

	anime, err := h.service.GetAnimeById(uint(animeId))
	if err != nil {
		handleError(c, err)
		return
	}
	c.JSON(http.StatusOK, anime)
}

func (h *HttpAnimeHandler) GetAnimeList(c *gin.Context) {
	animeQuery := dtos.AnimeQueryDTO{
		Seasonal: c.Query("seasonal"),
		Year:     c.Query("year"),
	}

	animes, err := h.service.GetAnimes(animeQuery)
	if err != nil {
		handleError(c, err)
		return
	}
	c.JSON(http.StatusOK, animes)

}

func (h *HttpAnimeHandler) UpdateAnime(c *gin.Context) {
	animeId, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		handleError(c, errs.NewBadRequestError(err.Error()))
		return
	}

	animeUpdate := new(entities.Anime)
	if err := c.BindJSON(animeUpdate); err != nil {
		handleError(c, errs.NewBadRequestError(err.Error()))
		return
	}
	animeUpdate.ID = uint(animeId)
	err = h.service.UpdateAnime(*animeUpdate)
	if err != nil {
		handleError(c, errs.NewBadRequestError(err.Error()))
		return
	}
	c.IndentedJSON(http.StatusOK, gin.H{"message": "Update anime success"})
}

func (h *HttpAnimeHandler) DeleteAnime(c *gin.Context) {
	animeId, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		handleError(c, errs.NewBadRequestError(err.Error()))
		return
	}

	if err = h.service.DeleteAnime(uint(animeId)); err != nil {
		handleError(c, err)
		return
	}
	c.IndentedJSON(http.StatusOK, gin.H{"message": "Delete anime success"})
}

func (h *HttpAnimeHandler) GetAnimeByUserId(c *gin.Context) {
	user_id, err := strconv.Atoi(c.Param("user_id"))
	if err != nil {
		handleError(c, errs.NewBadRequestError(err.Error()))
		return
	}
	animes, err := h.service.GetAnimeByUserId(uint(user_id))
	if err != nil {
		handleError(c, err)
		return
	}
	c.JSON(http.StatusOK, animes)
}

func (h *HttpAnimeHandler) GetAnimeByCategory(c *gin.Context) {
	category_id, err := strconv.Atoi(c.Param("category_id"))
	if err != nil {
		handleError(c, errs.NewBadRequestError(err.Error()))
		return
	}
	animes, err := h.service.GetAnimeByCategoryId(uint(category_id))
	if err != nil {
		handleError(c, err)
		return
	}
	c.JSON(http.StatusOK, animes)
}

func (h *HttpAnimeHandler) AddCategoryToAnime(c *gin.Context) {
	animeId, err := strconv.Atoi(c.Param("anime_id"))
	if err != nil {
		handleError(c, errs.NewBadRequestError(err.Error()))
		return
	}

	categoryRequest := new(dtos.AddCategoryToAnimeRequest)
	if err := c.BindJSON(categoryRequest); err != nil {
		handleError(c, errs.NewBadRequestError(err.Error()))
		return
	}

	categoryRequest.AnimeID = uint(animeId)
	if err := h.service.AddCategoryToAnime(*categoryRequest); err != nil {
		handleError(c, err)
		return
	}
	c.IndentedJSON(http.StatusOK, gin.H{"message": "Add anime to category success."})
}
