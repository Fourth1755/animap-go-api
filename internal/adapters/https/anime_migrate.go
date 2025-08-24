package adapters

import (
	"net/http"

	"github.com/Fourth1755/animap-go-api/internal/core/dtos"
	"github.com/Fourth1755/animap-go-api/internal/core/services"
	"github.com/Fourth1755/animap-go-api/internal/errs"
	"github.com/gin-gonic/gin"
)

type HttpAnimeMigrateHandler struct {
	service services.AnimeMigrateService
}

func NewHttpAnimeMigrateHandler(service services.AnimeMigrateService) *HttpAnimeMigrateHandler {
	return &HttpAnimeMigrateHandler{service: service}
}

func (h *HttpAnimeMigrateHandler) MigrateAnime(c *gin.Context) {
	request := new(dtos.MigrateAnimeRequest)
	if err := c.BindJSON(request); err != nil {
		handleError(c, errs.NewBadRequestError(err.Error()))
		return
	}
	err := h.service.MigrateAnime(*request)
	if err != nil {
		handleError(c, err)
		return
	}
	c.IndentedJSON(http.StatusOK, gin.H{"message": "Migrate anime success."})
}
