package adapters

import (
	"net/http"

	"github.com/Fourth1755/animap-go-api/internal/core/entities"
	"github.com/Fourth1755/animap-go-api/internal/core/services"
	"github.com/Fourth1755/animap-go-api/internal/errs"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type HttpCategoryUniverseHandler struct {
	service services.CategoryUniverseService
}

func NewHttpCategoryUniverseHandler(service services.CategoryUniverseService) *HttpCategoryUniverseHandler {
	return &HttpCategoryUniverseHandler{service: service}
}

func (h *HttpCategoryUniverseHandler) CreateCategoryUniverse(c *gin.Context) {
	category := new(entities.CategoryUniverse)
	if err := c.BindJSON(&category); err != nil {
		handleError(c, errs.NewBadRequestError(err.Error()))
		return
	}
	if err := h.service.CreateCategoryUniverse(category); err != nil {
		handleError(c, err)
		return
	}
	c.IndentedJSON(http.StatusCreated, gin.H{"message": "Create category success"})
}

func (h *HttpCategoryUniverseHandler) Getcategorise(c *gin.Context) {
	category, err := h.service.GetCategoryUniverses()
	if err != nil {
		handleError(c, errs.NewBadRequestError(err.Error()))
		return
	}
	c.JSON(http.StatusOK, category)
}

func (h *HttpCategoryUniverseHandler) GetCategoryUniverseById(c *gin.Context) {
	categoryId := c.Param("category_id")
	category, err := h.service.GetCategoryUniverseById(uuid.MustParse(categoryId))
	if err != nil {
		handleError(c, err)
		return
	}
	c.JSON(http.StatusOK, category)
}
