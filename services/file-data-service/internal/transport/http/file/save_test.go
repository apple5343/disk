package file

import (
	"bytes"
	"context"
	"data/internal/models"
	"data/internal/service"
	"data/internal/service/folder"
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

func randomFile() *models.FileMetadata {
	folder := "/" + gofakeit.Word() + "/" + gofakeit.Word()
	fileName := gofakeit.Word() + ".jpg"
	fullPath := folder + "/" + fileName
	return &models.FileMetadata{
		UserID:      gofakeit.UUID(),
		StoragePath: "users" + fullPath,
		FolderID:    gofakeit.UUID(),
		FileName:    fileName,
		FullPath:    fullPath,
		Bucket:      "test-bucket",
		Size:        int64(gofakeit.Uint16()),
		ContentType: "image/jpeg",
		Tags: map[string]string{
			"tag1": "value1",
			"tag2": "value2",
		},
		CreatedAt: time.Now(),
	}
}

func TestSave(t *testing.T) {
	t.Parallel()

	type storageServiceMockFunc func() service.FileService
	file := randomFile()
	serviceFile := *file
	serviceFile.ID = gofakeit.UUID()
	serviceFile.UpdatedAt = file.CreatedAt.Add(time.Minute)

	type args struct {
		file fileBody
	}

	req := fileBody{
		UserID:      file.UserID,
		StoragePath: file.StoragePath,
		FolderID:    file.FolderID,
		FileName:    file.FileName,
		FullPath:    file.FullPath,
		Bucket:      file.Bucket,
		Size:        file.Size,
		ContentType: file.ContentType,
		Tags:        file.Tags,
		CreatedAt:   file.CreatedAt,
	}

	resp := req
	resp.ID = serviceFile.ID
	resp.UpdatedAt = serviceFile.UpdatedAt

	tests := []struct {
		name        string
		args        args
		code        int
		want        fileBody
		wantErr     bool
		serviceFunc storageServiceMockFunc
	}{
		{
			name: "success",
			args: args{
				file: req,
			},
			code: http.StatusOK,
			want: resp,
			serviceFunc: func() service.FileService {
				m := mocks.NewFileService(t)
				m.EXPECT().
					SaveFile(mock.Anything, mock.Anything).
					Return(&serviceFile, nil).
					Run(func(_ context.Context, metadata *models.FileMetadata) {
						require.Equal(t, req.UserID, metadata.UserID)
						require.Equal(t, req.ID, metadata.ID)
					})
				return m
			},
		},
		{
			name: "folder not found",
			args: args{
				file: req,
			},
			code: http.StatusBadRequest,
			serviceFunc: func() service.FileService {
				m := mocks.NewFileService(t)
				m.EXPECT().
					SaveFile(mock.Anything, mock.Anything).
					Return(nil, folder.ErrFolderNotFound).
					Run(func(_ context.Context, metadata *models.FileMetadata) {
						require.Equal(t, req.FolderID, metadata.FolderID)
					})
				return m
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			reqBody, err := json.Marshal(tt.args.file)
			require.NoError(t, err)
			req := httptest.NewRequest(http.MethodPost, "/files", bytes.NewBuffer(reqBody))
			req.Header.Set("Content-Type", "application/json")
			rec := httptest.NewRecorder()

			c := echo.New().NewContext(req, rec)
			handler := NewHandler(tt.serviceFunc())
			middlewares.ErrorMiddleware()(handler.SaveFileMetadata())(c)

			require.Equal(t, tt.code, rec.Code)
			if !tt.wantErr {
				resp = fileBody{}
				err = json.NewDecoder(rec.Body).Decode(&resp)
				require.NoError(t, err)
				require.Equal(t, tt.want.ID, resp.ID)
			}
		})
	}
}
