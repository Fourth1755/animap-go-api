package middleware

import (
	"fmt"
	"os"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt"
)

func AuthRequired(c *fiber.Ctx) error {
	cookie := c.Cookies("jwt")
	jwtSecretKey := os.Getenv("JWT_SECRET")
	token, err := jwt.ParseWithClaims(cookie, jwt.MapClaims{}, func(t *jwt.Token) (interface{}, error) {
		return []byte(jwtSecretKey), nil
	})
	if err != nil || !token.Valid {
		return c.SendStatus(fiber.StatusUnauthorized)
	}
	claim := token.Claims.(jwt.MapClaims)
	fmt.Println(claim["role"])
	return c.Next()
}
