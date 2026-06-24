package app

import (
	httpclient "storage/internal/infrastructure/file-data/http"
	"storage/internal/infrastructure/health"
	"storage/internal/infrastructure/kafka/kafkago"
	"storage/internal/infrastructure/minio"
	"storage/internal/infrastructure/redis"

	"go.uber.org/fx"
)

func InfrastructureModule() fx.Option {
	return fx.Module(
		"infrastructure",
		health.NewModule(),
		redis.NewModule(),
		kafkago.NewModule(),
		httpclient.NewModule(),
		minio.NewModule(),
	)
}
