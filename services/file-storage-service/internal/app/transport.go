package app

import (
	httpserver "storage/internal/transport/http"
	"storage/internal/transport/redis"

	"go.uber.org/fx"
)

func TransportModule() fx.Option {
	return fx.Module(
		"transport",
		redis.NewModule(),
		httpserver.NewModule(),
	)
}
