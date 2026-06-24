package httpserver

import (
	"context"
	"data/internal/config"
	"data/internal/infrastructure/health"
	"data/internal/service"
	"data/internal/transport/http/file"
	"data/internal/transport/http/folder"
	"data/internal/transport/http/middlewares"
	"data/pkg/logger"
	"net/http"

	"github.com/labstack/echo/v4"
)

type Server struct {
	cfg           *config.HTTPConfig
	e             *echo.Echo
	file          *file.Handler
	folder        *folder.Handler
	l             logger.Logger
	jwtCfg        *config.JWTConfig
	healthChecker *health.Checker
}

func NewServer(
	cfg *config.HTTPConfig,
	fileService service.FileService,
	folderService service.FolderService,
	collectorService service.CollectorService,
	l logger.Logger,
	jwtCfg *config.JWTConfig,
	checker *health.Checker,
) *Server {
	return &Server{
		cfg:           cfg,
		e:             echo.New(),
		file:          file.NewHandler(fileService),
		folder:        folder.NewHandler(folderService, collectorService),
		l:             l,
		jwtCfg:        jwtCfg,
		healthChecker: checker,
	}
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

	files := s.e.Group("/files")
	{
		files.GET("/:id", s.file.GetFileMetadata())
		files.DELETE("/:id", s.file.DeleteFileMetadata())
		files.GET("", s.file.SearchFiles())
	}

	folders := s.e.Group("/folders")
	{
		folders.GET("/tree/:id", s.folder.GetFolderTree())
		folders.GET("/:id", s.folder.GetFolder())
		folders.DELETE("/:id", s.folder.DeleteFolder())
		folders.GET("", s.folder.GetRootFolder())
		folders.POST("", s.folder.SaveFolder())
	}
}

func (s *Server) Run() error {
	s.routes()
	return s.e.Start(s.cfg.Addr)
}

func (s *Server) Stop(ctx context.Context) error {
	return s.e.Shutdown(ctx)
}
