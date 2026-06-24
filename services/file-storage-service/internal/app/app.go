package app

import (
	"time"

	"go.uber.org/fx"
)

const (
	StopTimeout = time.Hour
)

func NewApp() *fx.App {
	return fx.New(
		ConfigModule(),
		LoggerModule(),
		InfrastructureModule(),
		RepositoryModule(),
		ServiceModule(),
		TransportModule(),
		fx.StopTimeout(StopTimeout),
	)
}
