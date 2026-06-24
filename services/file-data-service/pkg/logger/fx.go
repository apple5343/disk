package logger

import (
	"data/internal/config"

	"go.uber.org/fx"
)

func NewModule() fx.Option {
	return fx.Module(
		"zap",
		fx.Provide(
			config.NewLoggerConfig,
			NewLogger,
		),
	)
}
