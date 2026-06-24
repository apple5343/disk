package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"gitlab.crja72.ru/golang/2025/autumn/projects/go36/disk/services/auth-service/internal/client"
	"gitlab.crja72.ru/golang/2025/autumn/projects/go36/disk/services/auth-service/internal/config"
	"gitlab.crja72.ru/golang/2025/autumn/projects/go36/disk/services/auth-service/internal/handler"
	"gitlab.crja72.ru/golang/2025/autumn/projects/go36/disk/services/auth-service/internal/service"
)

const defaultDownTimeout = 10 * time.Second

func main() {
	cfg, err := config.Read()
	if err != nil {
		log.Fatalf("main: %v", err)
	}

	userClient := client.NewHTTPUserClient(cfg.UserSvcURL, cfg.InternalAPIKey)
	authService := service.NewAuthService(userClient, cfg.JWTSecret)
	authHandler := handler.NewAuthHandler(authService)

	app := fiber.New(fiber.Config{
		AppName: "Auth Service",
	})
	app.Use(logger.New())

	app.Get("/health", func(c *fiber.Ctx) error {
		return c.SendStatus(http.StatusOK)
	})
	app.Post("/api/v1/auth/register", authHandler.Register)
	app.Post("/api/v1/auth/login", authHandler.Login)

	go func() {
		if err = app.Listen(":" + cfg.AuthSvcPort); err != nil {
			log.Printf("Server error: %v", err)
		}
	}()

	log.Printf("Auth Service started on port %s", cfg.AuthSvcPort)

	// Graceful shutdown
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	log.Println("Shutting down server...")

	ctx, cancel := context.WithTimeout(context.Background(), defaultDownTimeout)
	defer cancel()

	if err = app.ShutdownWithContext(ctx); err != nil {
		log.Printf("Graceful shutdown failed: %v", err)
	} else {
		log.Println("Server stopped gracefully")
	}
}
