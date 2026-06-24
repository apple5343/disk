package file

import (
	"data/internal/models"
	httpModels "data/internal/transport/http/models"
	"net/http"

	"github.com/apple5343/errorx"
	"github.com/labstack/echo/v4"
)

func (h *Handler) GetFileMetadata() echo.HandlerFunc {
	return func(c echo.Context) error {
		id := c.Param("id")
		if id == "" {
			return errorx.NewError("id is required", errorx.BadRequest)
		}

		file, err := h.s.GetFileByID(c.Request().Context(), id)
		if err != nil {
			return err
		}
		return c.JSON(http.StatusOK, httpModels.FileMetadataToHTTP(file))
	}
}

func (h *Handler) SearchFiles() echo.HandlerFunc {
	return func(c echo.Context) error {
		fileName := c.QueryParam("file")
		folderID := c.QueryParam("folder")
		path := c.QueryParam("path")

		files, err := h.s.SearchFiles(c.Request().Context(), &models.SearchRequest{
			FileName: fileName,
			FolderID: folderID,
			Path:     path,
		})
		if err != nil {
			return err
		}
		resp := make([]*httpModels.FileMetadata, len(files))
		for i, file := range files {
			resp[i] = httpModels.FileMetadataToHTTP(file)
		}
		return c.JSON(http.StatusOK, resp)
	}
}
