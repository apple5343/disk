package file

import "go.uber.org/fx"

func NewModule() fx.Option {
	return fx.Module(
		"file-adapter",
		fx.Provide(
			NewAdapter,
		),
	)
}
