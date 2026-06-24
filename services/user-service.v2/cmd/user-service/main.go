package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"gitlab.crja72.ru/golang/2025/autumn/projects/go36/disk/services/user-service.v2/internal/config"
	"gitlab.crja72.ru/golang/2025/autumn/projects/go36/disk/services/user-service.v2/internal/handler"
	"gitlab.crja72.ru/golang/2025/autumn/projects/go36/disk/services/user-service.v2/internal/infrastructure/db"
	"gitlab.crja72.ru/golang/2025/autumn/projects/go36/disk/services/user-service.v2/internal/repository"
	"gitlab.crja72.ru/golang/2025/autumn/projects/go36/disk/services/user-service.v2/internal/service"
)

const shutdownTimeout = 10 * time.Second

func main() {
	cfg, err := config.Read()
	if err != nil {
		log.Fatalf("main: %v", err)
	}

	sqlDB, err := db.Connect(cfg.DatabaseURL)
	if err != nil {
		log.Fatalf("main: %v", err)
	}
	defer sqlDB.Close()

	userRepo := repository.NewPostgresRepository(sqlDB)
	userService := service.NewUserService(userRepo)

	app := fiber.New()

	app.Use(logger.New())

	handler.RegisterHandlers(app, cfg, userService)

	go func() {
		if err = app.Listen(":" + cfg.UserSvcPort); err != nil {
			log.Printf("Server error: %v", err)
		}
	}()

	log.Printf("User Service started on port %s", cfg.UserSvcPort)

	// Graceful shutdown
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	log.Println("Shutting down server...")

	ctx, cancel := context.WithTimeout(context.Background(), shutdownTimeout)
	defer cancel()

	if err = app.ShutdownWithContext(ctx); err != nil {
		log.Printf("Graceful shutdown failed: %v", err)
	} else {
		log.Println("Server stopped gracefully")
	}
}
