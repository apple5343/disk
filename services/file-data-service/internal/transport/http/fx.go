package httpserver

import (
	"context"
	"data/internal/config"
	"log/slog"

	"go.uber.org/fx"
)

func NewModule() fx.Option {
	return fx.Module(
		"http",
		fx.Provide(
			config.NewHTTPConfig,
			NewServer,
		),
		fx.Invoke(
			func(lc fx.Lifecycle, s *Server) {
				lc.Append(fx.Hook{
					OnStart: func(_ context.Context) error {
						go func() {
							if err := s.Run(); err != nil {
								slog.Error(err.Error())
							}
						}()
						return nil
					},
					OnStop: func(ctx context.Context) error {
						return s.Stop(ctx)
					},
				})
			},
		),
	)
}
