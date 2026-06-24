package middleware

import (
	jwtware "github.com/gofiber/contrib/jwt"
	"github.com/gofiber/fiber/v2"
)

const BadTokenMsg = "Expired or invalid token"

func JWTProtected(secretKey string) fiber.Handler {
	return jwtware.New(jwtware.Config{
		SigningKey:   jwtware.SigningKey{Key: []byte(secretKey)},
		ErrorHandler: jwtError,
	})
}

func jwtError(c *fiber.Ctx, _ error) error {
	return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
		"error": true,
		"msg":   BadTokenMsg,
	})
}
