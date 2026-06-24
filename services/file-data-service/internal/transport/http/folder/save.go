package folder

import (
	"net/http"

	"github.com/apple5343/errorx"
	"github.com/labstack/echo/v4"

	httpModels "data/internal/transport/http/models"
)

func (h *Handler) SaveFolder() echo.HandlerFunc {
	return func(c echo.Context) error {
		var folder httpModels.Folder
		if err := c.Bind(&folder); err != nil {
			return errorx.NewError(err.Error(), errorx.BadRequest)
		}
		file, err := h.s.SaveFolder(c.Request().Context(), httpModels.FolderFromHTTP(&folder))
		if err != nil {
			return err
		}
		return c.JSON(http.StatusOK, httpModels.FolderToHTTP(file))
	}
}
