package config

import (
	"fmt"

	"github.com/ilyakaznacheev/cleanenv"
)

type PostgresConfig struct {
	DSN string `env:"POSTGRES_DSN" env-default:"postgres://postgres:postgres@localhost:5432/postgres?sslmode=disable"`
}

func NewPostgresConfig() (*PostgresConfig, error) {
	cfg := &PostgresConfig{}
	if err := cleanenv.ReadEnv(cfg); err != nil {
		return nil, fmt.Errorf("read config: %w", err)
	}
	return cfg, nil
}
