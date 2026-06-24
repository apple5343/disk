package file

import (
	"net/http"

	httpModels "data/internal/transport/http/models"

	"github.com/apple5343/errorx"
	"github.com/labstack/echo/v4"
)

func (h *Handler) SaveFileMetadata() echo.HandlerFunc {
	return func(c echo.Context) error {
		var metadata httpModels.FileMetadata
		if err := c.Bind(&metadata); err != nil {
			return errorx.NewError(err.Error(), errorx.BadRequest)
		}
		file, err := h.s.SaveFile(c.Request().Context(), httpModels.FileMetadataFromHTTP(&metadata))
		if err != nil {
			return err
		}
		return c.JSON(http.StatusOK, httpModels.FileMetadataToHTTP(file))
	}
}
