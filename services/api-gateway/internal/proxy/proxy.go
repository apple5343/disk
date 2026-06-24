package proxy

import (
	"github.com/gofiber/fiber/v2"
	"gitlab.crja72.ru/golang/2025/autumn/projects/go36/disk/services/api-gateway/internal/config"
	"gitlab.crja72.ru/golang/2025/autumn/projects/go36/disk/services/api-gateway/internal/middleware"
)

const (
	APIPrefix  = "/api/"
	APIVersion = "v1"
)

// RegisterRoutes регистрирует все маршруты.
func RegisterRoutes(app *fiber.App, cfg *config.Config) {
	// Health check
	app.Get("/health", func(c *fiber.Ctx) error {
		return c.SendStatus(fiber.StatusOK)
	})

	// -jwt
	authProxy := newFullpathProxy(cfg.AuthSvcURL, "auth-service")
	app.Post(APIPrefix+APIVersion+"/auth/register", authProxy.Handler())
	app.Post(APIPrefix+APIVersion+"/auth/login", authProxy.Handler())

	// +jwt
	api := app.Group(APIPrefix + APIVersion)
	api.Use(middleware.JWTProtected(cfg.JWTSecret))

	fileDataProxy := newStrippedProxy(cfg.FileDataSvcURL, "file-data-service")
	api.All("/files", fileDataProxy.Handler())
	api.All("/files/*", fileDataProxy.Handler())
	api.All("/folders", fileDataProxy.Handler())
	api.All("/folders/*", fileDataProxy.Handler())

	fileStorageProxy := newStrippedProxy(cfg.FileStorageSvcURL, "file-storage-service")
	api.All("/file", fileStorageProxy.Handler())
	api.All("/file/*", fileStorageProxy.Handler())

	app.Get("/", func(c *fiber.Ctx) error {
		return c.Redirect(APIPrefix+APIVersion+"/files", fiber.StatusTemporaryRedirect)
	})
}
