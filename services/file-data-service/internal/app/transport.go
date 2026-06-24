package app

import (
	httpserver "data/internal/transport/http"
	"data/internal/transport/kafka"

	"go.uber.org/fx"
)

func TransportModule() fx.Option {
	return fx.Module(
		"transport",
		kafka.NewModule(),
		httpserver.NewModule(),
	)
}
