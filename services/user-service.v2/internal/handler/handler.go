package handler

import (
	"github.com/gofiber/fiber/v2"
	"gitlab.crja72.ru/golang/2025/autumn/projects/go36/disk/services/user-service.v2/internal/config"
	"gitlab.crja72.ru/golang/2025/autumn/projects/go36/disk/services/user-service.v2/internal/middleware"
	"gitlab.crja72.ru/golang/2025/autumn/projects/go36/disk/services/user-service.v2/internal/service"
)

func RegisterHandlers(app *fiber.App, cfg *config.Config, userService *service.UserService) {
	app.Get("/health", func(c *fiber.Ctx) error {
		return c.SendStatus(fiber.StatusOK)
	})

	userHandler := NewUserHandler(userService)
	app.Get("/api/v1/users/:id", userHandler.GetUser)

	// Внутренние маршруты
	internal := app.Group("/internal")
	internal.Use(middleware.InternalAuth(cfg))
	internal.Post("/users", userHandler.CreateUser)
	internal.Get("/users/by-login/:login", userHandler.GetUserByLogin)
}
