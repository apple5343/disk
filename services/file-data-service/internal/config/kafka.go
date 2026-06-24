package config

import "github.com/ilyakaznacheev/cleanenv"

type KafkaConsumerConfig struct {
	Brokers []string `env:"KAFKA_BROKERS"            env-default:"localhost:29092"`
	Topic   string   `env:"KAFKA_TOPIC"              env-default:"uploading_files"`
	GroupID string   `env:"FILE_DATA_KAFKA_GROUP_ID" env-default:"file-data-service"`
}

func NewKafkaConsumerConfig() (*KafkaConsumerConfig, error) {
	cfg := &KafkaConsumerConfig{}

	if err := cleanenv.ReadEnv(cfg); err != nil {
		return nil, err
	}
	return cfg, nil
}
