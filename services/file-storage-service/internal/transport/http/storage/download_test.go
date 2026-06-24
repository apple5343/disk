package storage

import (
	"bytes"
	"context"
	"io"
	"net/http"
	"net/http/httptest"
	"storage/internal/service"
	"storage/internal/service/mocks"
	"storage/internal/transport/http/middlewares"
	"strconv"
	"testing"

	storageService "storage/internal/service/storage"

	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

func TestDownload(t *testing.T) {
	t.Parallel()
	type storageServiceMockFunc func() service.StorageService

	type args struct {
		fileID string
	}

	type want struct {
		contentType        string
		contentDisposition string
		contentLength      int
		fileID             string
		body               []byte
	}

	fileData := []byte{1, 2, 3}
	expectedFile := randomFile()
	expectedFile.Size = int64(len(fileData))

	tests := []struct {
		name        string
		args        args
		want        want
		statusCode  int
		wantErr     bool
		serviceMock storageServiceMockFunc
	}{
		{
			name: "success",
			args: args{
				fileID: expectedFile.ID,
			},
			want: want{
				contentType:        expectedFile.ContentType,
				contentLength:      int(expectedFile.Size),
				contentDisposition: "attachment; filename=\"" + expectedFile.FileName + "\"",
				fileID:             expectedFile.ID,
				body:               fileData,
			},
			statusCode: http.StatusOK,
			serviceMock: func() service.StorageService {
				m := mocks.NewStorageService(t)
				m.EXPECT().
					DownloadFile(mock.Anything, mock.Anything).
					Return(expectedFile, bytes.NewBuffer(fileData), nil).
					Run(func(_ context.Context, id string) {
						require.Equal(t, expectedFile.ID, id)
					})
				return m
			},
		},
		{
			name: "file not found",
			args: args{
				fileID: expectedFile.ID,
			},
			statusCode: http.StatusBadRequest,
			serviceMock: func() service.StorageService {
				m := mocks.NewStorageService(t)
				m.EXPECT().
					DownloadFile(mock.Anything, mock.Anything).
					Return(nil, nil, storageService.ErrFileNotFound)
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

			req := httptest.NewRequest(http.MethodGet, "/file/"+tt.args.fileID, nil)
			recorder := httptest.NewRecorder()
			c := e.NewContext(req, recorder)
			c.SetParamNames("id")
			c.SetParamValues(tt.args.fileID)
			middlewares.ErrorMiddleware()(handler.DownloadHandler())(c)
			require.Equal(t, tt.statusCode, recorder.Code)

			if !tt.wantErr {
				contentType := recorder.Header().Get("Content-Type")
				require.Equal(t, tt.want.contentType, contentType)
				contentDisposition := recorder.Header().Get("Content-Disposition")
				require.Equal(t, tt.want.contentDisposition, contentDisposition)
				contentLength := recorder.Header().Get("Content-Length")
				require.Equal(t, strconv.Itoa(tt.want.contentLength), contentLength)
				fileID := recorder.Header().Get("X-File-Id")
				require.Equal(t, tt.want.fileID, fileID)
				got, err := io.ReadAll(recorder.Body)
				require.NoError(t, err)
				require.Equal(t, tt.want.body, got)
			}
		})
	}
}
