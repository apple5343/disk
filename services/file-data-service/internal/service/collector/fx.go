package collector

import "go.uber.org/fx"

func NewModule() fx.Option {
	return fx.Module(
		"collector",
		fx.Provide(
			NewService,
		),
	)
}
