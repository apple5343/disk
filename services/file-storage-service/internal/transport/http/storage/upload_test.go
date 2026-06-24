package storage

import (
	"bytes"
	"context"
	"encoding/json"
	"io"
	"mime/multipart"
	"net/http"
	"net/http/httptest"
	"storage/internal/models"
	"storage/internal/service"
	"storage/internal/service/mocks"
	storageService "storage/internal/service/storage"
	"storage/internal/transport/http/middlewares"
	"testing"
	"time"

	"github.com/brianvoe/gofakeit"
	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

func createMultipartRequest(t *testing.T, folderID, filename string, data []byte) (io.Reader, string) {
	var buf bytes.Buffer
	w := multipart.NewWriter(&buf)

	require.NoError(t, w.WriteField("folder_id", folderID))

	part, err := w.CreateFormFile("file", filename)
	require.NoError(t, err)

	_, err = io.Copy(part, bytes.NewReader(data))
	require.NoError(t, err)

	require.NoError(t, w.Close())
	return &buf, w.FormDataContentType()
}

func randomFile() *models.FileMetadata {
	folder := "/" + gofakeit.Word() + "/" + gofakeit.Word()
	fileName := gofakeit.Word() + ".jpg"
	fullPath := folder + "/" + fileName
	return &models.FileMetadata{
		ID:          gofakeit.UUID(),
		UserID:      gofakeit.UUID(),
		StoragePath: "users/" + fullPath,
		FolderID:    gofakeit.UUID(),
		FileName:    fileName,
		FullPath:    fullPath,
		Bucket:      "test-bucket",
		Size:        gofakeit.Int64(),
		ContentType: "image/jpeg",
		Tags: map[string]string{
			"tag1": "value1",
			"tag2": "value2",
		},
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}
}

func TestUpload(t *testing.T) {
	t.Parallel()
	type storageServiceMockFunc func() service.StorageService

	type args struct {
		fileData []byte
		fileName string
		folderID string
		userID   string
	}

	type respBody struct {
		Key string `json:"key"`
	}

	fileData := []byte{1, 2, 3}
	expectedFile := randomFile()
	expectedFile.Size = int64(len(fileData))

	tests := []struct {
		name        string
		args        args
		want        respBody
		statusCode  int
		wantErr     bool
		serviceMock storageServiceMockFunc
	}{
		{
			name: "success",
			args: args{
				fileData: fileData,
				folderID: expectedFile.FolderID,
				fileName: expectedFile.FileName,
			},
			want: respBody{
				Key: expectedFile.ID,
			},
			statusCode: http.StatusOK,
			serviceMock: func() service.StorageService {
				m := mocks.NewStorageService(t)
				m.EXPECT().
					UploadFile(mock.Anything, mock.Anything, mock.Anything).
					Return(expectedFile, nil).
					Run(func(_ context.Context, meta *models.FileMetadata, data io.Reader) {
						require.Equal(t, expectedFile.FileName, meta.FileName)
						require.Equal(t, expectedFile.FolderID, meta.FolderID)
						require.Equal(t, int64(len(fileData)), meta.Size)

						got, err := io.ReadAll(data)
						require.NoError(t, err)
						require.Equal(t, fileData, got)
					})
				return m
			},
		},
		{
			name: "folder not found",
			args: args{
				fileData: fileData,
				folderID: expectedFile.FolderID,
				fileName: expectedFile.FileName,
				userID:   expectedFile.UserID,
			},
			statusCode: http.StatusBadRequest,
			serviceMock: func() service.StorageService {
				m := mocks.NewStorageService(t)
				m.EXPECT().
					UploadFile(mock.Anything, mock.Anything, mock.Anything).
					Return((*models.FileMetadata)(nil), storageService.ErrFolderNotFound)
				return m
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			e := echo.New()
			service := tt.serviceMock()
			handler := NewHandler(service)

			body, contentType := createMultipartRequest(t, tt.args.folderID, tt.args.fileName, tt.args.fileData)
			req := httptest.NewRequest(http.MethodPost, "/file", body)
			req.Header.Set("Content-Type", contentType)

			recorder := httptest.NewRecorder()
			c := e.NewContext(req, recorder)
			middlewares.ErrorMiddleware()(handler.UploadHanler())(c)

			require.Equal(t, tt.statusCode, recorder.Code)
			if !tt.wantErr {
				var got respBody
				require.NoError(t, json.Unmarshal(recorder.Body.Bytes(), &got))
				require.Equal(t, tt.want, got)
			}
		})
	}
}
