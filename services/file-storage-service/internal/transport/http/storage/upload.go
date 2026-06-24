package storage

import (
	"net/http"
	"storage/internal/models"

	"github.com/apple5343/errorx"
	"github.com/labstack/echo/v4"
)

func (h *Handler) UploadHanler() echo.HandlerFunc {
	return func(c echo.Context) error {
		formFile, err := c.FormFile("file")
		if err != nil {
			return errorx.NewError("file is required", errorx.BadRequest)
		}

		src, err := formFile.Open()
		if err != nil {
			return err
		}

		file := &models.FileMetadata{
			Size:        formFile.Size,
			FileName:    formFile.Filename,
			FolderID:    c.FormValue("folder_id"),
			ContentType: formFile.Header.Get("Content-Type"),
		}

		file, err = h.s.UploadFile(c.Request().Context(), file, src)
		if err != nil {
			return err
		}

		return c.JSON(http.StatusOK, map[string]string{
			"key": file.ID,
		})
	}
}

func (h *Handler) CancelUpload() echo.HandlerFunc {
	return func(c echo.Context) error {
		id := c.Param("id")
		if id == "" {
			return errorx.NewError("id is required", errorx.BadRequest)
		}
		err := h.s.CancelUpload(c.Request().Context(), id)
		if err != nil {
			return err
		}
		return c.NoContent(http.StatusNoContent)
	}
}
