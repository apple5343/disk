package folder

import "go.uber.org/fx"

func NewModule() fx.Option {
	return fx.Module(
		"folder",
		fx.Provide(
			NewService,
		),
	)
}
