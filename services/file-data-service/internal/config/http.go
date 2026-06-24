package config

import (
	"fmt"

	"github.com/ilyakaznacheev/cleanenv"
)

type HTTPConfig struct {
	Addr string `env:"FILE_DATA_HTTP_ADDR" env-default:":8003"`
}

func NewHTTPConfig() (*HTTPConfig, error) {
	cfg := &HTTPConfig{}
	if err := cleanenv.ReadEnv(cfg); err != nil {
		return nil, fmt.Errorf("read config: %w", err)
	}
	return cfg, nil
}
