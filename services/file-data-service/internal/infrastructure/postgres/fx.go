package postgres

import (
	"context"
	"data/internal/config"

	"github.com/jmoiron/sqlx"
	"go.uber.org/fx"
)

func NewModule() fx.Option {
	return fx.Module(
		"postgres",
		fx.Provide(
			config.NewPostgresConfig,
			NewClient,
		),
		fx.Invoke(
			func(lc fx.Lifecycle, client *sqlx.DB) {
				lc.Append(fx.Hook{
					OnStop: func(_ context.Context) error {
						return client.Close()
					},
				})
			},
		),
	)
}
