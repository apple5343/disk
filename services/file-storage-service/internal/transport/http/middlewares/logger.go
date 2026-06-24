package middlewares

import (
	"storage/pkg/logger"

	"github.com/labstack/echo/v4"
)

func LoggerMiddleware(l logger.Logger) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			ctx := logger.ContextWithLogger(c.Request().Context(), l)
			c.SetRequest(c.Request().WithContext(ctx))
			return next(c)
		}
	}
}
