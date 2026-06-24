package config

import (
	"fmt"

	"github.com/ilyakaznacheev/cleanenv"
)

type FileDataClientConfig struct {
	BaseURL string `env:"FILE_DATA_SVC_URL" env-default:"http://localhost:8003"`
}

func NewFileDataClientConfig() (*FileDataClientConfig, error) {
	cfg := &FileDataClientConfig{}
	if err := cleanenv.ReadEnv(cfg); err != nil {
		return nil, fmt.Errorf("read config: %w", err)
	}
	return cfg, nil
}
