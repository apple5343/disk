package config

import (
	"fmt"

	"github.com/ilyakaznacheev/cleanenv"
)

type LoggerConfig struct {
	Level string `env:"LOGGER_LEVEL" env-default:"dev"`
}

func NewLoggerConfig() (*LoggerConfig, error) {
	cfg := &LoggerConfig{}
	if err := cleanenv.ReadEnv(cfg); err != nil {
		return nil, fmt.Errorf("read config: %w", err)
	}
	return cfg, nil
}
