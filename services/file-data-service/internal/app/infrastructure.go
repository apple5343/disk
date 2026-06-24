package app

import (
	"data/internal/infrastructure/health"
	"data/internal/infrastructure/kafka/kafkago"
	"data/internal/infrastructure/postgres"
	"data/internal/infrastructure/redis"

	"go.uber.org/fx"
)

func InfrastructureModule() fx.Option {
	return fx.Module(
		"infrastructure",
		kafkago.NewModule(),
		redis.NewModule(),
		postgres.NewModule(),
		health.NewModule(),
	)
}
