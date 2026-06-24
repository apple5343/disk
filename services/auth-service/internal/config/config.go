package config

import (
	"fmt"

	"github.com/ilyakaznacheev/cleanenv"
)

type Config struct {
	AuthSvcPort    string `env:"AUTH_SVC_PORT"    env-default:"8001"`
	UserSvcURL     string `env:"USER_SVC_URL"     env-default:"http://user-service:8002"`
	UserSvcPort    string `env:"USER_SVC_PORT"    env-default:"8002"`
	JWTSecret      string `env:"JWT_SECRET"       env-default:"super-secret-dev-only"`
	InternalAPIKey string `env:"INTERNAL_API_KEY" env-default:"secret-key-for-internal-communication"`
}

func Read() (*Config, error) {
	cfg := &Config{}
	if err := cleanenv.ReadEnv(cfg); err != nil {
		return nil, fmt.Errorf("read config: %w", err)
	}
	return cfg, nil
}
