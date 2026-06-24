package middleware

import (
	"github.com/gofiber/fiber/v2"
	"gitlab.crja72.ru/golang/2025/autumn/projects/go36/disk/services/user-service.v2/internal/config"
)

//nolint:gosec // Key name.
const InternalAPIKeyHeader = "X-Internal-API-Key"

func InternalAuth(cfg *config.Config) fiber.Handler {
	return func(c *fiber.Ctx) error {
		key := c.Get(InternalAPIKeyHeader)
		if key != cfg.InternalAPIKey {
			return c.Status(fiber.StatusForbidden).JSON(fiber.Map{
				"error": "forbidden",
			})
		}
		return c.Next()
	}
}
