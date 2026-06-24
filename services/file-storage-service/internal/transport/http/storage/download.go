package storage

import (
	"fmt"
	"io"
	"net/http"
	"strconv"

	"github.com/apple5343/errorx"
	"github.com/labstack/echo/v4"
)

func (h *Handler) DownloadHandler() echo.HandlerFunc {
	return func(c echo.Context) error {
		id := c.Param("id")
		if id == "" {
			return errorx.NewError("id is required", errorx.BadRequest)
		}

		metadata, reader, err := h.s.DownloadFile(c.Request().Context(), id)
		if err != nil {
			return err
		}

		closer, ok := reader.(io.Closer)
		if ok {
			defer closer.Close()
		}

		c.Response().Header().Set("Content-Type", metadata.ContentType)
		c.Response().Header().Set("Content-Disposition",
			fmt.Sprintf("attachment; filename=\"%s\"", metadata.FileName))
		c.Response().Header().Set("Content-Length", strconv.FormatInt(metadata.Size, 10))
		c.Response().Header().Set("X-File-Id", metadata.ID)

		return c.Stream(http.StatusOK, metadata.ContentType, reader)
	}
}
