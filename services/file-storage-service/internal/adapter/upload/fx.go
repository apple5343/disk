package upload

import "go.uber.org/fx"

func NewModule() fx.Option {
	return fx.Module(
		"upload",
		fx.Provide(
			NewUploadingAdapter,
		),
	)
}
