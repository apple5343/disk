package redis

import (
	"context"

	"go.uber.org/fx"
)

func NewModule() fx.Option {
	return fx.Module(
		"redis",
		fx.Provide(
			NewSubscriber,
		),
		fx.Invoke(
			func(lc fx.Lifecycle, s *Subscriber) {
				ctx, cancel := context.WithCancel(context.Background())
				lc.Append(fx.Hook{
					OnStart: func(_ context.Context) error {
						go func() {
							if err := s.Listen(ctx); err != nil {
								panic(err)
							}
						}()
						return nil
					},
					OnStop: func(_ context.Context) error {
						cancel()
						return nil
					},
				})
			},
		),
	)
}
