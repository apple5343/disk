package storage

import (
	"context"
	"storage/internal/service"

	"go.uber.org/fx"
)

func NewModule() fx.Option {
	return fx.Module(
		"storage",
		fx.Provide(
			NewService,
		),
		fx.Invoke(
			func(lc fx.Lifecycle, s service.StorageService) {
				lc.Append(fx.Hook{
					OnStop: func(_ context.Context) error {
						return s.Close(context.Background())
					},
				})
			},
		),
	)
}
