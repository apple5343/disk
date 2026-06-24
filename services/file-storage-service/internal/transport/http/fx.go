package httpserver

import (
	"context"
	"errors"
	"net/http"
	"storage/internal/config"

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
							if err := s.Run(); err != nil && !errors.Is(err, http.ErrServerClosed) {
								panic(err)
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
