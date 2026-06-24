package kafkago

import (
	"context"
	"storage/internal/config"

	infrastructure "storage/internal/infrastructure/kafka"

	"github.com/segmentio/kafka-go"
)

const (
	pingTopic = "ping"
)

type producer struct {
	writer *kafka.Writer
}

func NewProducer(cfg *config.KafkaProducerConfig) (infrastructure.Producer, error) {
	writer := &kafka.Writer{
		Addr:     kafka.TCP(cfg.Addr),
		Balancer: &kafka.LeastBytes{},
	}
	return &producer{
		writer: writer,
	}, nil
}

func (p *producer) Close(_ context.Context) error {
	return p.writer.Close()
}

func (p *producer) Publish(ctx context.Context, topic string, key []byte, value []byte) error {
	return p.writer.WriteMessages(ctx, kafka.Message{
		Topic: topic,
		Key:   key,
		Value: value,
	})
}

func (p *producer) Ping(ctx context.Context) error {
	return p.Publish(ctx, pingTopic, []byte("ping"), []byte("ping"))
}
