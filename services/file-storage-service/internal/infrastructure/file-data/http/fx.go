package httpclient

import (
	"storage/internal/config"

	"go.uber.org/fx"
)

func NewModule() fx.Option {
	return fx.Module(
		"file-data-client",
		fx.Provide(
			config.NewFileDataClientConfig,
			NewClient,
		),
	)
}
