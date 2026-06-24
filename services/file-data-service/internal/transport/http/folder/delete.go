package folder

import (
	"net/http"

	"github.com/apple5343/errorx"
	"github.com/labstack/echo/v4"
)

func (h *Handler) DeleteFolder() echo.HandlerFunc {
	return func(c echo.Context) error {
		id := c.Param("id")
		if id == "" {
			return errorx.NewError("id is required", errorx.BadRequest)
		}
		err := h.collector.DeleteFolder(c.Request().Context(), id)
		if err != nil {
			return err
		}
		return c.JSON(http.StatusNoContent, nil)
	}
}
