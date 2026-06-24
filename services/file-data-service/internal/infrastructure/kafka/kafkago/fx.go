package kafkago

import (
	"context"

	"data/internal/config"
	infrastructure "data/internal/infrastructure/kafka"

	"go.uber.org/fx"
)

func NewModule() fx.Option {
	return fx.Module(
		"kafka-consumer",
		fx.Provide(
			config.NewKafkaConsumerConfig,
			NewKafkaConsumer,
		),
		fx.Invoke(
			func(lc fx.Lifecycle, consumer infrastructure.Consumer) {
				lc.Append(fx.Hook{
					OnStop: func(ctx context.Context) error {
						return consumer.Close(ctx)
					},
				})
			},
		),
	)
}
