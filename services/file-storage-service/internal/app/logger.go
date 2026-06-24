package app

import (
	"storage/pkg/logger"

	"go.uber.org/fx"
)

func LoggerModule() fx.Option {
	return fx.Module(
		"logger",
		logger.NewModule(),
	)
}
