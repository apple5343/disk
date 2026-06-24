package config

import (
	"fmt"

	"github.com/ilyakaznacheev/cleanenv"
)

type JWTConfig struct {
	Secret string `env:"JWT_SECRET" env-default:"super-secret-jwt-key-for-dev-only"`
}

func NewJWTConfig() (*JWTConfig, error) {
	cfg := &JWTConfig{}
	if err := cleanenv.ReadEnv(cfg); err != nil {
		return nil, fmt.Errorf("read config: %w", err)
	}
	return cfg, nil
}
