package kafkago

import (
	"context"

	"storage/internal/config"
	infrastructure "storage/internal/infrastructure/kafka"

	"go.uber.org/fx"
)

func NewModule() fx.Option {
	return fx.Module(
		"kafkago-producer",
		fx.Provide(
			config.NewKafkaProducerConfig,
			NewProducer,
		),
		fx.Invoke(
			func(lc fx.Lifecycle, producer infrastructure.Producer) {
				lc.Append(fx.Hook{
					OnStop: func(ctx context.Context) error {
						return producer.Close(ctx)
					},
				})
			},
		),
	)
}
