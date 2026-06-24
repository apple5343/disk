package config

import (
	"fmt"

	"github.com/ilyakaznacheev/cleanenv"
)

type ConsumerConfig struct {
	RateLimitPerSecond int `env:"RATE_LIMIT_PER_SECOND" env-default:"10"`
	MaxInProcess       int `env:"MAX_IN_PROCESS"        env-default:"10"`
}

func NewConsumerConfig() (*ConsumerConfig, error) {
	cfg := &ConsumerConfig{}
	if err := cleanenv.ReadEnv(cfg); err != nil {
		return nil, fmt.Errorf("read config: %w", err)
	}
	return cfg, nil
}
