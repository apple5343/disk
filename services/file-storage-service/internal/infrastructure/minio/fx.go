package minio

import (
	"storage/internal/config"

	"go.uber.org/fx"
)

func NewModule() fx.Option {
	return fx.Module(
		"minio",
		fx.Provide(
			config.NewMinioConfig,
			NewClient,
		),
	)
}
