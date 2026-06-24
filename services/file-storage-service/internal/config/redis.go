package config

import (
	"fmt"

	"github.com/ilyakaznacheev/cleanenv"
)

type RedisConfig struct {
	Addr     string `env:"REDIS_ADDR"     env-default:"localhost:6379"`
	Password string `env:"REDIS_PASSWORD" env-default:"redisPassword"`
	DB       int    `env:"REDIS_DB"       env-default:"0"`
}

func NewRedisConfig() (*RedisConfig, error) {
	cfg := &RedisConfig{}
	if err := cleanenv.ReadEnv(cfg); err != nil {
		return nil, fmt.Errorf("read config: %w", err)
	}
	return cfg, nil
}
