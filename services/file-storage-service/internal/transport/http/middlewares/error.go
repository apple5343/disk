package middlewares

import (
	"errors"
	"fmt"
	"net/http"
	"storage/pkg/logger"

	"github.com/apple5343/errorx"
	"github.com/labstack/echo/v4"
)

func ErrorMiddleware() echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			err := next(c)
			if nil == err {
				return nil
			}

			if c.Response().Committed {
				return err
			}
			l, okk := logger.FromContext(c.Request().Context())
			if okk {
				l.Error(c.Request().Context(), err.Error())
			}
			httpErr := &echo.HTTPError{}
			if errors.As(err, &httpErr) {
				return c.JSON(httpErr.Code, map[string]string{
					"status":  "error",
					"message": fmt.Sprintf("%v", httpErr.Message),
				})
			}

			if commonErr, ok := errorx.ToCommonError(err); ok {
				if commonErr.Code() == errorx.Internal {
					return c.JSON(
						http.StatusInternalServerError,
						map[string]string{"status": "error", "message": "Internal server error"},
					)
				}

				return c.JSON(
					parseCommonCode(commonErr.Code()),
					map[string]string{"status": "error", "message": commonErr.Error()},
				)
			}

			return c.JSON(
				http.StatusInternalServerError,
				map[string]string{"status": "error", "message": "Internal server error"},
			)
		}
	}
}

func parseCommonCode(code errorx.Code) int {
	switch code {
	case errorx.Unauthorized:
		return http.StatusUnauthorized
	case errorx.Forbidden:
		return http.StatusForbidden
	case errorx.Conflict:
		return http.StatusConflict
	case errorx.OK:
		return http.StatusOK
	case errorx.BadRequest:
		return http.StatusBadRequest
	case errorx.Internal:
		return http.StatusInternalServerError
	case errorx.NotFound:
		return http.StatusNotFound
	case errorx.Cancelled:
		return http.StatusGone
	case errorx.DeadlineExceeded:
		return http.StatusRequestTimeout
	default:
		return http.StatusInternalServerError
	}
}
