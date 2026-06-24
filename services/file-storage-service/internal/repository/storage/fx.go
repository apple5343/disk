package storage

import "go.uber.org/fx"

func NewModule() fx.Option {
	return fx.Module(
		"minio",
		fx.Provide(
			NewRepository,
		),
	)
}
