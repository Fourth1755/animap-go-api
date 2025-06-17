package adapters

import (
	"net/http"

	"github.com/Fourth1755/animap-go-api/internal/core/entities"
	"github.com/Fourth1755/animap-go-api/internal/core/services"
	"github.com/Fourth1755/animap-go-api/internal/errs"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type HttpArtistHandler struct {
	service services.ArtistService
}

func NewHttpArtistHandler(service services.ArtistService) *HttpArtistHandler {
	return &HttpArtistHandler{service: service}
}

func (h *HttpArtistHandler) CreateArtist(c *gin.Context) {
	var artist *entities.Artist
	if err := c.BindJSON(&artist); err != nil {
		handleError(c, errs.NewBadRequestError(err.Error()))
		return
	}
	if err := h.service.CreateArtist(artist); err != nil {
		handleError(c, err)
		return
	}
	c.IndentedJSON(http.StatusCreated, gin.H{"message": "Create artist success"})
}

func (h *HttpArtistHandler) GetArtistById(c *gin.Context) {
	artistId := c.Param("id")
	artist, err := h.service.GetArtistById(uuid.MustParse(artistId))
	if err != nil {
		handleError(c, err)
		return
	}
	c.JSON(http.StatusOK, artist)
}

func (h *HttpArtistHandler) GetArtistList(c *gin.Context) {
	artists, err := h.service.GetArtists()
	if err != nil {
		handleError(c, errs.NewBadRequestError(err.Error()))
		return
	}

	c.JSON(http.StatusOK, artists)
}
