package app

import (
	"time"

	"go.uber.org/fx"
)

const (
	StopTimeout = 30 * time.Second
)

func NewApp() *fx.App {
	return fx.New(
		ConfigModule(),
		LoggerModule(),
		InfrastructureModule(),
		RepositoryModule(),
		ServiceModule(),
		TransportModule(),
		fx.StartTimeout(StopTimeout),
	)
}
