package config

import (
	"fmt"

	"github.com/ilyakaznacheev/cleanenv"
)

type MinioConfig struct {
	Endpoint        string `env:"MINIO_ENDPOINT"          env-default:"localhost:9000"`
	AccessKeyID     string `env:"MINIO_ACCESS_KEY_ID"     env-default:"minioadmin"`
	SecretAccessKey string `env:"MINIO_SECRET_ACCESS_KEY" env-default:"minioadmin"`
	UseSSL          bool   `env:"MINIO_USE_SSL"           env-default:"false"`
}

func NewMinioConfig() (*MinioConfig, error) {
	cfg := &MinioConfig{}
	if err := cleanenv.ReadEnv(cfg); err != nil {
		return nil, fmt.Errorf("read config: %w", err)
	}
	return cfg, nil
}
