package adapters

import (
	"net/http"

	"github.com/Fourth1755/animap-go-api/internal/core/entities"
	"github.com/Fourth1755/animap-go-api/internal/core/services"
	"github.com/Fourth1755/animap-go-api/internal/errs"
	"github.com/gin-gonic/gin"
)

type HttpUserHandler struct {
	service services.UserService
}

func NewHttpUserHandler(service services.UserService) *HttpUserHandler {
	return &HttpUserHandler{service: service}
}

func (h *HttpUserHandler) CreateUser(c *gin.Context) {
	user := new(entities.User)
	if err := c.BindJSON(user); err != nil {
		handleError(c, errs.NewBadRequestError(err.Error()))
		return
	}
	err := h.service.CreateUser(user)
	if err != nil {
		handleError(c, err)
		return
	}
	c.IndentedJSON(http.StatusOK, gin.H{"message": "Register success."})
}

func (h *HttpUserHandler) Login(c *gin.Context) {
	user := new(entities.User)
	if err := c.BindJSON(user); err != nil {
		handleError(c, errs.NewBadRequestError(err.Error()))
		return
	}

	token, err := h.service.Login(user)
	if err != nil {
		handleError(c, errs.NewUnauthorizedError(err.Error()))
		return
	}
	c.SetCookie("jwt", token, 3600*24, "/", "localhost", false, true)
	c.IndentedJSON(http.StatusOK, gin.H{"message": "Login success.", "token": token})
}
