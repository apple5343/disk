package folder

import (
	httpModels "data/internal/transport/http/models"
	ctxutil "data/internal/utils/ctx"
	"net/http"

	"github.com/apple5343/errorx"
	"github.com/labstack/echo/v4"
)

func (h *Handler) GetFolder() echo.HandlerFunc {
	return func(c echo.Context) error {
		id := c.Param("id")
		if id == "" {
			return errorx.NewError("id is required", errorx.BadRequest)
		}
		file, err := h.s.GetFolderByID(c.Request().Context(), id)
		if err != nil {
			return err
		}
		return c.JSON(http.StatusOK, httpModels.FolderToHTTP(file))
	}
}

func (h *Handler) GetRootFolder() echo.HandlerFunc {
	return func(c echo.Context) error {
		userID := ctxutil.UserIDFromContext(c.Request().Context())
		if userID == "" {
			return errorx.NewError("user id is required", errorx.BadRequest)
		}
		file, err := h.s.RootFolder(c.Request().Context(), userID)
		if err != nil {
			return err
		}
		return c.JSON(http.StatusOK, httpModels.FolderToHTTP(file))
	}
}

func (h *Handler) GetFolderTree() echo.HandlerFunc {
	return func(c echo.Context) error {
		id := c.Param("id")
		if id == "" {
			return errorx.NewError("id is required", errorx.BadRequest)
		}
		tree, err := h.collector.GetFolderTree(c.Request().Context(), id)
		if err != nil {
			return err
		}
		return c.JSON(http.StatusOK, httpModels.FolderTreeToHTTP(tree))
	}
}
