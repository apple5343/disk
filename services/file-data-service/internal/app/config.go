package app

import (
	"data/internal/config"

	"go.uber.org/fx"
)

func ConfigModule() fx.Option {
	return fx.Module(
		"config",
		fx.Provide(
			config.NewJWTConfig,
		),
	)
}
