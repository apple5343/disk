package folder

import (
	"bytes"
	"context"
	"data/internal/models"
	folderService "data/internal/service/folder"
	"data/internal/service/mocks"
	"data/internal/transport/http/middlewares"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

func TestSave(t *testing.T) {
	t.Parallel()

	t.Run("success", func(t *testing.T) {
		t.Parallel()

		folder := randomFolder()

		fs := mocks.NewFolderService(t)
		cs := mocks.NewCollectorService(t)

		fs.EXPECT().
			SaveFolder(mock.Anything, mock.Anything).
			Return(folder, nil).
			Run(func(_ context.Context, f *models.Folder) {
				require.Equal(t, folder.Name, f.Name)
				require.Equal(t, *folder.ParentID, *f.ParentID)
			})

		reqFolder := folderBody{
			UserID:   folder.UserID,
			Name:     folder.Name,
			ParentID: *folder.ParentID,
		}
		respFolder, err := json.Marshal(reqFolder)
		require.NoError(t, err)
		req := httptest.NewRequest(http.MethodPost, "/folders", bytes.NewBuffer(respFolder))
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rec := httptest.NewRecorder()

		c := echo.New().NewContext(req, rec)
		handler := NewHandler(fs, cs)
		middlewares.ErrorMiddleware()(handler.SaveFolder())(c)

		require.Equal(t, http.StatusOK, rec.Code)
		var respFolderBody folderBody
		err = json.NewDecoder(rec.Body).Decode(&respFolderBody)
		require.NoError(t, err)
		require.Equal(t, folder.ID, respFolderBody.ID)
		require.Equal(t, folder.Name, respFolderBody.Name)
		require.Equal(t, *folder.ParentID, respFolderBody.ParentID)
		require.Equal(t, folder.UserID, respFolderBody.UserID)
		require.Equal(t, folder.PathDepth, respFolderBody.PathDepth)
		require.Equal(t, folder.FullPath, respFolderBody.FullPath)
	})

	t.Run("invalid parent", func(t *testing.T) {
		t.Parallel()

		folder := randomFolder()

		fs := mocks.NewFolderService(t)
		cs := mocks.NewCollectorService(t)

		fs.EXPECT().
			SaveFolder(mock.Anything, mock.Anything).
			Return(nil, folderService.ErrInvalidParent).
			Run(func(_ context.Context, f *models.Folder) {
				require.Equal(t, folder.Name, f.Name)
				require.Equal(t, folder.ParentID, f.ParentID)
			})

		reqFolder := folderBody{
			UserID:   folder.UserID,
			Name:     folder.Name,
			ParentID: *folder.ParentID,
		}
		respFolder, err := json.Marshal(reqFolder)
		require.NoError(t, err)
		req := httptest.NewRequest(http.MethodPost, "/folders", bytes.NewBuffer(respFolder))
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rec := httptest.NewRecorder()

		c := echo.New().NewContext(req, rec)
		handler := NewHandler(fs, cs)
		middlewares.ErrorMiddleware()(handler.SaveFolder())(c)

		require.Equal(t, http.StatusBadRequest, rec.Code)
	})
}
