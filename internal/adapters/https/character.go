package adapters

import (
	"net/http"

	"github.com/Fourth1755/animap-go-api/internal/core/dtos"
	"github.com/Fourth1755/animap-go-api/internal/core/services"
	"github.com/Fourth1755/animap-go-api/internal/errs"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type HttpCharacterHandler struct {
	service services.CharacterService
}

func NewHttpCharacterHandler(service services.CharacterService) *HttpCharacterHandler {
	return &HttpCharacterHandler{service: service}
}

func (h *HttpCharacterHandler) CreateCharacter(c *gin.Context) {
	var request dtos.CreateCharacterRequest
	if err := c.BindJSON(&request); err != nil {
		handleError(c, errs.NewBadRequestError(err.Error()))
		return
	}
	err := h.service.CreateCharacter(request)
	if err != nil {
		handleError(c, err)
		return
	}
	c.IndentedJSON(http.StatusCreated, gin.H{"message": "Create Character success"})
}

func (h *HttpCharacterHandler) GetCharacterByAnimeId(c *gin.Context) {
	animeId := c.Param("anime_id")
	character, err := h.service.GetCharacterByAnimeId(uuid.MustParse(animeId))
	if err != nil {
		handleError(c, err)
		return
	}
	c.JSON(http.StatusOK, character)
}
