package sqlx

import "go.uber.org/fx"

func NewModule() fx.Option {
	return fx.Module(
		"sqlx-folder",
		fx.Provide(
			NewRepository,
		),
	)
}
