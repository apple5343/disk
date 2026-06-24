package file

import "go.uber.org/fx"

func NewModule() fx.Option {
	return fx.Module(
		"file",
		fx.Provide(
			NewService,
		),
	)
}
