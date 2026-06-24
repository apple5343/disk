package redis

import (
	"go.uber.org/fx"
)

func NewModule() fx.Option {
	return fx.Module(
		"file-processing",
		fx.Provide(
			NewRepository,
		),
	)
}
