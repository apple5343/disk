package file

import (
	"context"
	"data/internal/service/file"
	"data/internal/service/mocks"
	"data/internal/transport/http/middlewares"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/brianvoe/gofakeit"
	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

func TestDele(t *testing.T) {
	t.Parallel()

	t.Run("success", func(t *testing.T) {
		t.Parallel()

		fileID := gofakeit.UUID()
		service := mocks.NewFileService(t)
		service.EXPECT().
			DeleteFile(mock.Anything, fileID).
			Return(nil).
			Run(func(_ context.Context, id string) {
				require.Equal(t, fileID, id)
			})

		req := httptest.NewRequest(http.MethodDelete, "/files/"+fileID, nil)
		rec := httptest.NewRecorder()

		c := echo.New().NewContext(req, rec)
		c.SetParamNames("id")
		c.SetParamValues(fileID)

		handler := NewHandler(service)
		middlewares.ErrorMiddleware()(handler.DeleteFileMetadata())(c)

		require.Equal(t, http.StatusNoContent, rec.Code)
	})

	t.Run("empty id", func(t *testing.T) {
		t.Parallel()

		service := mocks.NewFileService(t)

		req := httptest.NewRequest(http.MethodDelete, "/files/", nil)
		rec := httptest.NewRecorder()

		c := echo.New().NewContext(req, rec)

		handler := NewHandler(service)
		middlewares.ErrorMiddleware()(handler.DeleteFileMetadata())(c)

		require.Equal(t, http.StatusBadRequest, rec.Code)
	})

	t.Run("file not found", func(t *testing.T) {
		t.Parallel()

		fileID := gofakeit.UUID()
		service := mocks.NewFileService(t)
		service.EXPECT().
			DeleteFile(mock.Anything, fileID).
			Return(file.ErrFileNotFound).
			Run(func(_ context.Context, id string) {
				require.Equal(t, fileID, id)
			})

		req := httptest.NewRequest(http.MethodDelete, "/files/"+fileID, nil)
		rec := httptest.NewRecorder()

		c := echo.New().NewContext(req, rec)
		c.SetParamNames("id")
		c.SetParamValues(fileID)

		handler := NewHandler(service)
		middlewares.ErrorMiddleware()(handler.DeleteFileMetadata())(c)

		require.Equal(t, http.StatusBadRequest, rec.Code)
	})
}
