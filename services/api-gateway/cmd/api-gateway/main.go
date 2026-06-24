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
	"gitlab.crja72.ru/golang/2025/autumn/projects/go36/disk/services/api-gateway/internal/config"
	"gitlab.crja72.ru/golang/2025/autumn/projects/go36/disk/services/api-gateway/internal/proxy"
)

const defaultDownTimeout = 5 * time.Second

type App struct {
	config *config.Config
}

func NewApp(cfg *config.Config) *App {
	return &App{config: cfg}
}

func (a *App) run(ctx context.Context, downTimeout time.Duration) error {
	app := fiber.New()

	app.Use(logger.New())

	proxy.RegisterRoutes(app, a.config)

	go func() {
		log.Printf("Starting server on port %s", a.config.Port)
		if err := app.Listen(":" + a.config.Port); err != nil {
			log.Printf("Server failed to start: %v", err)
		}
	}()

	<-ctx.Done()

	log.Println("Shutting down server...")

	shutdownCtx, cancel := context.WithTimeout(context.Background(), downTimeout)
	defer cancel()

	if err := app.ShutdownWithContext(shutdownCtx); err != nil {
		log.Printf("Server shutdown failed: %v", err)
		return err
	}

	log.Println("Server shutdown successfully")
	return nil
}

// testableMain - версия main для тестирования.
// Натягивание совы на глобус для повышения coverage.
func testableMain(
	ctx context.Context,
	downTimeout time.Duration,
	readConfig func() (*config.Config, error),
	newApp func(*config.Config) *App,
) error {
	cfg, err := readConfig()
	if err != nil {
		return err
	}
	app := newApp(cfg)
	return app.run(ctx, downTimeout)
}

func main() {
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()
	if err := testableMain(ctx, defaultDownTimeout, config.Read, NewApp); err != nil {
		log.Printf("Application failed: %v", err)
	}
}
