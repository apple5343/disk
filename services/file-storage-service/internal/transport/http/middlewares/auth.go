package middlewares

import (
	"storage/internal/config"
	ctxutil "storage/internal/utils/ctx"
	"storage/pkg/jwt"
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
			token := strings.TrimPrefix(authHeader, bearerPrefix)
			userClaims, err := jwt.VerifyToken(token, []byte(jwtConfig.Secret))
			if err != nil {
				return errorx.NewError("invalid token", errorx.Unauthorized)
			}
			ctx := ctxutil.ContextWithUserID(c.Request().Context(), userClaims.ID)
			ctx = ctxutil.ContextWithToken(ctx, token)
			c.SetRequest(c.Request().WithContext(ctx))
			return next(c)
		}
	}
}
