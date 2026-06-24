package config

import "github.com/ilyakaznacheev/cleanenv"

type KafkaProducerConfig struct {
	Addr string `env:"KAFKA_BROKERS" env-default:"localhost:29092"`
}

func NewKafkaProducerConfig() (*KafkaProducerConfig, error) {
	cfg := &KafkaProducerConfig{}

	if err := cleanenv.ReadEnv(cfg); err != nil {
		return nil, err
	}
	return cfg, nil
}
