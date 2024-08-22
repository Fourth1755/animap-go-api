package adapters

import (
	"time"

	"github.com/Fourth1755/animap-go-api/internal/core/entities"
	"github.com/Fourth1755/animap-go-api/internal/core/services"
	"github.com/Fourth1755/animap-go-api/internal/errs"
	"github.com/gofiber/fiber/v2"
)

type HttpUserHandler struct {
	service services.UserService
}

func NewHttpUserHandler(service services.UserService) *HttpUserHandler {
	return &HttpUserHandler{service: service}
}

func (h *HttpUserHandler) CreateUser(c *fiber.Ctx) error {
	user := new(entities.User)
	if err := c.BodyParser(user); err != nil {
		return handleError(c, errs.NewBadRequestError(err.Error()))
	}
	err := h.service.CreateUser(user)
	if err != nil {
		return handleError(c, err)
	}
	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "Register success",
	})
}

func (h *HttpUserHandler) Login(c *fiber.Ctx) error {
	user := new(entities.User)
	if err := c.BodyParser(user); err != nil {
		return handleError(c, errs.NewBadRequestError(err.Error()))
	}

	token, err := h.service.Login(user)
	if err != nil {
		return handleError(c, errs.NewUnauthorizedError(err.Error()))
	}
	c.Cookie(&fiber.Cookie{
		Name:     "jwt",
		Value:    token,
		Expires:  time.Now().Add(time.Hour * 72),
		HTTPOnly: true,
	})
	return c.JSON(fiber.Map{
		"message": "Login success",
		"token":   token,
	})
}
