package adapters

import (
	"net/http"
	"strconv"

	"github.com/Fourth1755/animap-go-api/internal/core/dtos"
	"github.com/Fourth1755/animap-go-api/internal/core/services"
	"github.com/Fourth1755/animap-go-api/internal/errs"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type HttpCommentHandler struct {
	service services.CommentService
}

func NewHttpCommentHandler(service services.CommentService) *HttpCommentHandler {
	return &HttpCommentHandler{service: service}
}

func (h *HttpCommentHandler) CreateComment(c *gin.Context) {
	var req dtos.CreateCommentAnimeRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		handleError(c, errs.NewBadRequestError(err.Error()))
		return
	}

	err := h.service.CreateComment(c, req)
	if err != nil {
		handleError(c, err)
		return
	}

	c.IndentedJSON(http.StatusCreated, gin.H{"message": "Create comment success."})
}

func (h *HttpCommentHandler) GetComments(c *gin.Context) {
	// Parse anime ID from path
	animeIDStr := c.Param("id")
	animeID, err := uuid.Parse(animeIDStr)
	if err != nil {
		handleError(c, errs.NewBadRequestError("Invalid anime ID format"))
		return
	}

	// Parse query parameters
	commentType := c.Query("type")

	pageStr := c.DefaultQuery("page", "1")
	page, err := strconv.Atoi(pageStr)
	if err != nil || page < 1 {
		page = 1
	}

	limitStr := c.DefaultQuery("limit", "10")
	limit, err := strconv.Atoi(limitStr)
	if err != nil || limit < 1 {
		limit = 10
	}

	// Call service
	comments, err := h.service.GetComments(animeID, commentType, page, limit)
	if err != nil {
		handleError(c, err)
		return
	}

	c.JSON(http.StatusOK, comments)
}
