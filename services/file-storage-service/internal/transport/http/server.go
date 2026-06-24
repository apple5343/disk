package httpserver

import (
	"context"
	"net/http"
	"storage/internal/config"
	"storage/internal/infrastructure/health"
	"storage/internal/service"
	"storage/internal/transport/http/middlewares"
	"storage/internal/transport/http/storage"
	"storage/pkg/logger"

	"github.com/labstack/echo/v4"
)

type Server struct {
	cfg           *config.HTTPConfig
	e             *echo.Echo
	storage       *storage.Handler
	jwtCfg        *config.JWTConfig
	l             logger.Logger
	healthChecker *health.Checker
}

func NewServer(
	cfg *config.HTTPConfig,
	storageService service.StorageService,
	l logger.Logger,
	jwtCfg *config.JWTConfig,
	checker *health.Checker,
) *Server {
	return &Server{
		cfg:           cfg,
		e:             echo.New(),
		storage:       storage.NewHandler(storageService),
		l:             l,
		jwtCfg:        jwtCfg,
		healthChecker: checker,
	}
}

func (s *Server) Run() error {
	s.routes()
	return s.e.Start(s.cfg.Addr)
}

func (s *Server) Stop(ctx context.Context) error {
	return s.e.Shutdown(ctx)
}

func (s *Server) routes() {
	s.e.GET("/health", func(c echo.Context) error {
		err := s.healthChecker.Check(c.Request().Context())
		if err != nil {
			return c.JSON(http.StatusServiceUnavailable, map[string]string{"status": "error", "message": err.Error()})
		}
		return c.JSON(http.StatusOK, map[string]string{"status": "ok"})
	})
	s.e.Use(middlewares.LoggerMiddleware(s.l))
	s.e.Use(middlewares.ErrorMiddleware())
	s.e.Use(middlewares.AuthMiddleware(s.jwtCfg))
	s.e.POST("/file", s.storage.UploadHanler())
	s.e.DELETE("/terminate/:id", s.storage.CancelUpload())
	s.e.GET("/file/:id", s.storage.DownloadHandler())
	s.e.DELETE("/file/:id", s.storage.DeleteFileHandler())
}
