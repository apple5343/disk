package file

import (
	"context"
	"data/internal/models"
	"data/internal/service/file"
	"data/internal/service/mocks"
	"data/internal/transport/http/middlewares"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/brianvoe/gofakeit"
	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

type fileBody struct {
	ID          string            `json:"id"`
	UserID      string            `json:"user_id"`
	StoragePath string            `json:"storage_path"`
	FolderID    string            `json:"folder_id"`
	FileName    string            `json:"file_name"`
	FullPath    string            `json:"full_path"`
	Bucket      string            `json:"bucket"`
	Size        int64             `json:"size"`
	ContentType string            `json:"content_type"`
	Tags        map[string]string `json:"tags"`
	CreatedAt   time.Time         `json:"created_at"`
	UpdatedAt   time.Time         `json:"updated_at"`
}

func TestGet(t *testing.T) {
	t.Parallel()

	t.Run("success", func(t *testing.T) {
		t.Parallel()
		file := randomFile()
		file.ID = gofakeit.UUID()
		service := mocks.NewFileService(t)

		service.EXPECT().
			GetFileByID(mock.Anything, mock.Anything).
			Return(file, nil).
			Run(func(_ context.Context, id string) {
				require.Equal(t, file.ID, id)
			})

		req := httptest.NewRequest(http.MethodGet, "/files/"+file.ID, nil)
		rec := httptest.NewRecorder()

		c := echo.New().NewContext(req, rec)
		c.SetParamNames("id")
		c.SetParamValues(file.ID)

		handler := NewHandler(service)
		middlewares.ErrorMiddleware()(handler.GetFileMetadata())(c)

		require.Equal(t, http.StatusOK, rec.Code)
		resp := fileBody{}
		err := json.NewDecoder(rec.Body).Decode(&resp)
		require.NoError(t, err)
		require.Equal(t, file.ID, resp.ID)
		require.Equal(t, file.UserID, resp.UserID)
		require.Equal(t, file.StoragePath, resp.StoragePath)
		require.Equal(t, file.FolderID, resp.FolderID)
		require.Equal(t, file.FileName, resp.FileName)
		require.Equal(t, file.FullPath, resp.FullPath)
		require.Equal(t, file.Bucket, resp.Bucket)
		require.Equal(t, file.Size, resp.Size)
		require.Equal(t, file.ContentType, resp.ContentType)
		require.Equal(t, file.Tags, resp.Tags)
	})

	t.Run("empty id", func(t *testing.T) {
		t.Parallel()
		service := mocks.NewFileService(t)

		req := httptest.NewRequest(http.MethodGet, "/files/", nil)
		rec := httptest.NewRecorder()

		c := echo.New().NewContext(req, rec)

		handler := NewHandler(service)
		middlewares.ErrorMiddleware()(handler.GetFileMetadata())(c)

		require.Equal(t, http.StatusBadRequest, rec.Code)
	})

	t.Run("file not found", func(t *testing.T) {
		t.Parallel()
		fileID := gofakeit.UUID()

		service := mocks.NewFileService(t)
		service.EXPECT().
			GetFileByID(mock.Anything, mock.Anything).
			Return(nil, file.ErrFileNotFound).
			Run(func(_ context.Context, id string) {
				require.Equal(t, fileID, id)
			})

		req := httptest.NewRequest(http.MethodGet, "/files/"+fileID, nil)
		rec := httptest.NewRecorder()

		c := echo.New().NewContext(req, rec)
		c.SetParamNames("id")
		c.SetParamValues(fileID)

		handler := NewHandler(service)
		middlewares.ErrorMiddleware()(handler.GetFileMetadata())(c)

		require.Equal(t, http.StatusBadRequest, rec.Code)
	})
}

func TestSearch(t *testing.T) {
	t.Parallel()

	t.Run("success", func(t *testing.T) {
		t.Parallel()
		file := randomFile()
		file.ID = gofakeit.UUID()
		service := mocks.NewFileService(t)
		service.EXPECT().
			SearchFiles(mock.Anything, mock.Anything).
			Return([]*models.FileMetadata{file}, nil)

		req := httptest.NewRequest(
			http.MethodGet,
			"/files?path="+file.FullPath+"&folder="+file.FolderID+"&file="+file.FileName,
			nil,
		)
		rec := httptest.NewRecorder()

		c := echo.New().NewContext(req, rec)

		handler := NewHandler(service)
		middlewares.ErrorMiddleware()(handler.SearchFiles())(c)

		require.Equal(t, http.StatusOK, rec.Code)
		var respFiles []fileBody
		err := json.NewDecoder(rec.Body).Decode(&respFiles)
		resp := respFiles[0]
		require.NoError(t, err)
		require.Equal(t, file.ID, resp.ID)
		require.Equal(t, file.UserID, resp.UserID)
		require.Equal(t, file.StoragePath, resp.StoragePath)
		require.Equal(t, file.FolderID, resp.FolderID)
		require.Equal(t, file.FileName, resp.FileName)
		require.Equal(t, file.FullPath, resp.FullPath)
		require.Equal(t, file.Bucket, resp.Bucket)
		require.Equal(t, file.Size, resp.Size)
		require.Equal(t, file.ContentType, resp.ContentType)
		require.Equal(t, file.Tags, resp.Tags)
	})

	t.Run("invalid folder id", func(t *testing.T) {
		t.Parallel()

		service := mocks.NewFileService(t)
		fileID := gofakeit.UUID()
		service.EXPECT().
			SearchFiles(mock.Anything, mock.Anything).
			Return(nil, file.ErrInvalidID).
			Run(func(_ context.Context, req *models.SearchRequest) {
				require.Equal(t, fileID, req.FolderID)
			})

		req := httptest.NewRequest(http.MethodGet, "/files?folder="+fileID, nil)
		rec := httptest.NewRecorder()

		c := echo.New().NewContext(req, rec)

		handler := NewHandler(service)
		middlewares.ErrorMiddleware()(handler.SearchFiles())(c)

		require.Equal(t, http.StatusBadRequest, rec.Code)
	})
}
