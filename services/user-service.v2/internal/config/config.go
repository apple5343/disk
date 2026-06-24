package config

import (
	"fmt"

	"github.com/ilyakaznacheev/cleanenv"
)

type Config struct {
	UserSvcAddr    string `env:"USER_SVC_ADDR"    env-default:"user-service:8002"`
	UserSvcPort    string `env:"USER_SVC_PORT"    env-default:"8002"`
	DatabaseURL    string `env:"USER_DB_DSN"      env-default:"postgres://postgres:postgres@postgres:5432/user_db?sslmode=disable"`
	InternalAPIKey string `env:"INTERNAL_API_KEY" env-default:"secret-key-for-internal-communication"`
}

func Read() (*Config, error) {
	cfg := &Config{}
	if err := cleanenv.ReadEnv(cfg); err != nil {
		return nil, fmt.Errorf("read config: %w", err)
	}
	return cfg, nil
}
