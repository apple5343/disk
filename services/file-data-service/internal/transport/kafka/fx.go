package kafka

import (
	"context"
	"data/internal/config"
	"data/pkg/logger"
	"errors"

	"go.uber.org/fx"
)

func NewModule() fx.Option {
	return fx.Module(
		"kafka-consumer",
		fx.Provide(
			config.NewConsumerConfig,
			NewConsumer,
		),
		fx.Invoke(
			func(lc fx.Lifecycle, c *Consumer, l logger.Logger) {
				ctx := logger.ContextWithLogger(context.Background(), l)
				ctx, cancel := context.WithCancel(ctx)
				lc.Append(
					fx.Hook{
						OnStart: func(_ context.Context) error {
							go func() {
								if err := c.Consume(ctx); err != nil {
									if errors.Is(err, context.Canceled) {
										return
									}
									panic(err)
								}
							}()
							return nil
						},
						OnStop: func(ctx context.Context) error {
							cancel()
							return c.Close(ctx)
						},
					},
				)
			},
		),
	)
}
