package config

import (
	"fmt"

	"github.com/ilyakaznacheev/cleanenv"
)

type Config struct {
	Port              string `env:"GW_PORT"              env-default:"8000"`
	AuthSvcURL        string `env:"AUTH_SVC_URL"         env-default:"http://auth-service:8001"`
	FileDataSvcURL    string `env:"FILE_DATA_SVC_URL"    env-default:"http://file-data-service:8003"`
	FileStorageSvcURL string `env:"FILE_STORAGE_SVC_URL" env-default:"http://file-storage-service:8004"`
	JWTSecret         string `env:"JWT_SECRET"           env-default:"super-secret-jwt-key-for-dev-only"`
}

func Read() (*Config, error) {
	cfg := &Config{}
	if err := cleanenv.ReadEnv(cfg); err != nil {
		return nil, fmt.Errorf("read config: %w", err)
	}
	return cfg, nil
}
