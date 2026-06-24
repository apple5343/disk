package main

import (
	"errors"
	"fmt"
	"log"

	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"gitlab.crja72.ru/golang/2025/autumn/projects/go36/disk/services/user-service.v2/internal/config"
)

const migrationsDir = "file://migrations"

func RunMigrations(dbURL string) error {
	m, err := migrate.New(migrationsDir, dbURL)
	if err != nil {
		return fmt.Errorf("create migrate instance: %w", err)
	}
	defer m.Close()

	if err = m.Up(); err != nil && !errors.Is(err, migrate.ErrNoChange) {
		return fmt.Errorf("run migrations: %w", err)
	}
	return nil
}

func main() {
	cfg, err := config.Read()
	if err != nil {
		log.Fatalf("main: %v", err)
	}

	err = RunMigrations(cfg.DatabaseURL)
	if err != nil {
		log.Fatalf("main: %v", err)
	}
	log.Println("Migrations applied successfully")
}
