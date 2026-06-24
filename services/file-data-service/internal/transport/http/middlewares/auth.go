package middlewares

import (
	"data/internal/config"
	ctxutil "data/internal/utils/ctx"
	"data/pkg/jwt"
	"strings"

	"github.com/apple5343/errorx"
	"github.com/labstack/echo/v4"
)

const (
	bearerPrefix = "Bearer "
)

func AuthMiddleware(jwtConfig *config.JWTConfig) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			if c.Path() == "/health" {
				return next(c)
			}
			authHeader := c.Request().Header.Get("Authorization")
			if authHeader == "" {
				return errorx.NewError("token is empty", errorx.Unauthorized)
			}
			if !strings.HasPrefix(authHeader, bearerPrefix) {
				return errorx.NewError("invalid token", errorx.Unauthorized)
			}
			userClaims, err := jwt.VerifyToken(strings.TrimPrefix(authHeader, bearerPrefix), []byte(jwtConfig.Secret))
			if err != nil {
				return errorx.NewError("invalid token", errorx.Unauthorized)
			}
			ctx := ctxutil.ContextWithUserID(c.Request().Context(), userClaims.ID)
			c.SetRequest(c.Request().WithContext(ctx))
			return next(c)
		}
	}
}

func AuthMiddlewareV2() echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			ctx := ctxutil.ContextWithUserID(c.Request().Context(), c.Request().Header.Get("User_id")) // TODO исправить
			c.SetRequest(c.Request().WithContext(ctx))
			return next(c)
		}
	}
}
