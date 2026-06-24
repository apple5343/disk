package kafkago

import (
	"context"

	infrastructure "data/internal/infrastructure/kafka"

	"github.com/segmentio/kafka-go"

	"data/internal/config"
)

type consumer struct {
	reader *kafka.Reader
	cfg    *config.KafkaConsumerConfig
}

func NewKafkaConsumer(cfg *config.KafkaConsumerConfig) (infrastructure.Consumer, error) {
	reader := kafka.NewReader(kafka.ReaderConfig{
		Brokers:     cfg.Brokers,
		Topic:       cfg.Topic,
		StartOffset: kafka.FirstOffset,
		GroupID:     cfg.GroupID,
	})

	return &consumer{
		reader: reader,
		cfg:    cfg,
	}, nil
}

func (c *consumer) Read(ctx context.Context) (infrastructure.Message, error) {
	msg, err := c.reader.FetchMessage(ctx)
	if err != nil {
		return nil, err
	}
	return c.NewMessage(&msg), nil
}

func (c *consumer) Ping(ctx context.Context) error {
	client := &kafka.Client{
		Addr: kafka.TCP(c.cfg.Brokers[0]),
	}
	_, err := client.Metadata(ctx, &kafka.MetadataRequest{
		Topics: []string{c.cfg.Topic},
	})
	return err
}

func (c *consumer) Close(_ context.Context) error {
	return c.reader.Close()
}
