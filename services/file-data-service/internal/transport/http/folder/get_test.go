package folder

import (
	"context"
	"data/internal/models"
	"data/internal/service/folder"
	"data/internal/service/mocks"
	"data/internal/transport/http/middlewares"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/brianvoe/gofakeit"
	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

func randomFolder() *models.Folder {
	folderName := gofakeit.Word()
	parentID := gofakeit.UUID()
	return &models.Folder{
		ID:        gofakeit.UUID(),
		UserID:    gofakeit.UUID(),
		ParentID:  &parentID,
		Name:      folderName,
		FullPath:  "/" + gofakeit.Word() + "/" + folderName,
		PathDepth: 2,
	}
}

func randomFile() *models.FileMetadata {
	return &models.FileMetadata{
		ID:       gofakeit.UUID(),
		FolderID: gofakeit.UUID(),
	}
}

type folderBody struct {
	ID        string `json:"id"`
	UserID    string `json:"user_id"`
	Name      string `json:"name"`
	FullPath  string `json:"full_path"`
	ParentID  string `json:"parent_id"`
	PathDepth int    `json:"path_depth"`
}

type fileBody struct {
	ID     string `json:"id"`
	UserID string `json:"user_id"`
}

func TestGet(t *testing.T) {
	t.Parallel()

	t.Run("success", func(t *testing.T) {
		t.Parallel()

		folder := randomFolder()
		fs := mocks.NewFolderService(t)
		cs := mocks.NewCollectorService(t)

		fs.EXPECT().
			GetFolderByID(mock.Anything, mock.Anything).
			Return(folder, nil).
			Run(func(_ context.Context, id string) {
				require.Equal(t, folder.ID, id)
			})

		req := httptest.NewRequest(http.MethodGet, "/folders/"+folder.ID, nil)
		rec := httptest.NewRecorder()

		c := echo.New().NewContext(req, rec)
		c.SetParamNames("id")
		c.SetParamValues(folder.ID)

		handler := NewHandler(fs, cs)
		middlewares.ErrorMiddleware()(handler.GetFolder())(c)

		require.Equal(t, http.StatusOK, rec.Code)

		var resp folderBody
		err := json.NewDecoder(rec.Body).Decode(&resp)
		require.NoError(t, err)
		require.Equal(t, folder.ID, resp.ID)
		require.Equal(t, folder.UserID, resp.UserID)
		require.Equal(t, folder.Name, resp.Name)
		require.Equal(t, folder.FullPath, resp.FullPath)
		require.Equal(t, folder.PathDepth, resp.PathDepth)
	})

	t.Run("empty id", func(t *testing.T) {
		t.Parallel()

		fs := mocks.NewFolderService(t)
		cs := mocks.NewCollectorService(t)

		req := httptest.NewRequest(http.MethodGet, "/folders/", nil)
		rec := httptest.NewRecorder()

		c := echo.New().NewContext(req, rec)

		handler := NewHandler(fs, cs)
		middlewares.ErrorMiddleware()(handler.GetFolder())(c)

		require.Equal(t, http.StatusBadRequest, rec.Code)
	})

	t.Run("folder not found", func(t *testing.T) {
		t.Parallel()

		fs := mocks.NewFolderService(t)
		cs := mocks.NewCollectorService(t)

		folderID := gofakeit.UUID()
		fs.EXPECT().
			GetFolderByID(mock.Anything, mock.Anything).
			Return(nil, folder.ErrFolderNotFound).
			Run(func(_ context.Context, id string) {
				require.Equal(t, folderID, id)
			})

		req := httptest.NewRequest(http.MethodGet, "/folders/"+folderID, nil)
		rec := httptest.NewRecorder()

		c := echo.New().NewContext(req, rec)
		c.SetParamNames("id")
		c.SetParamValues(folderID)

		handler := NewHandler(fs, cs)
		middlewares.ErrorMiddleware()(handler.GetFolder())(c)

		require.Equal(t, http.StatusBadRequest, rec.Code)
	})
}

func TestFolderTree(t *testing.T) {
	t.Parallel()

	t.Run("success", func(t *testing.T) {
		t.Parallel()

		folder := randomFolder()
		file := randomFile()
		parent := gofakeit.UUID()
		fs := mocks.NewFolderService(t)
		cs := mocks.NewCollectorService(t)

		cs.EXPECT().
			GetFolderTree(mock.Anything, mock.Anything).
			Return(&models.FolderTree{Childs: []*models.Folder{folder}, Files: []*models.FileMetadata{file}}, nil).
			Run(func(_ context.Context, parentID string) {
				require.Equal(t, parent, parentID)
			})

		req := httptest.NewRequest(http.MethodGet, "/folders/tree/"+parent, nil)
		rec := httptest.NewRecorder()

		c := echo.New().NewContext(req, rec)
		c.SetParamNames("id")
		c.SetParamValues(parent)

		handler := NewHandler(fs, cs)
		middlewares.ErrorMiddleware()(handler.GetFolderTree())(c)

		require.Equal(t, http.StatusOK, rec.Code)

		type treeBody struct {
			Files  []fileBody   `json:"files"`
			Childs []folderBody `json:"folders"`
		}

		resp := treeBody{}
		err := json.NewDecoder(rec.Body).Decode(&resp)
		require.NoError(t, err)
		require.Equal(t, folder.ID, resp.Childs[0].ID)
		require.Equal(t, file.ID, resp.Files[0].ID)
	})

	t.Run("empty id", func(t *testing.T) {
		t.Parallel()

		fs := mocks.NewFolderService(t)
		cs := mocks.NewCollectorService(t)

		req := httptest.NewRequest(http.MethodGet, "/folders/tree/", nil)
		rec := httptest.NewRecorder()

		c := echo.New().NewContext(req, rec)

		handler := NewHandler(fs, cs)
		middlewares.ErrorMiddleware()(handler.GetFolderTree())(c)

		require.Equal(t, http.StatusBadRequest, rec.Code)
	})

	t.Run("folder not found", func(t *testing.T) {
		t.Parallel()

		folderID := gofakeit.UUID()
		fs := mocks.NewFolderService(t)
		cs := mocks.NewCollectorService(t)

		cs.EXPECT().
			GetFolderTree(mock.Anything, mock.Anything).
			Return(nil, folder.ErrFolderNotFound).
			Run(func(_ context.Context, parentID string) {
				require.Equal(t, folderID, parentID)
			})

		req := httptest.NewRequest(http.MethodGet, "/folders/tree/"+folderID, nil)
		rec := httptest.NewRecorder()

		c := echo.New().NewContext(req, rec)
		c.SetParamNames("id")
		c.SetParamValues(folderID)

		handler := NewHandler(fs, cs)
		middlewares.ErrorMiddleware()(handler.GetFolderTree())(c)

		require.Equal(t, http.StatusBadRequest, rec.Code)
	})
}
