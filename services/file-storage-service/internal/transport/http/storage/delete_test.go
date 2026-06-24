package storage

import (
	"net/http"
	"net/http/httptest"
	"storage/internal/service"
	"storage/internal/service/mocks"
	"storage/internal/service/storage"
	"storage/internal/transport/http/middlewares"
	"testing"

	"github.com/brianvoe/gofakeit"
	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

func TestDelete(t *testing.T) {
	t.Parallel()

	type storageServiceMockFunc func() service.StorageService

	type args struct {
		fileID string
	}
	fileID := gofakeit.UUID()

	tests := []struct {
		name        string
		args        args
		statusCode  int
		serviceMock storageServiceMockFunc
	}{
		{
			name: "success",
			args: args{
				fileID: fileID,
			},
			statusCode: http.StatusNoContent,
			serviceMock: func() service.StorageService {
				m := mocks.NewStorageService(t)
				m.EXPECT().
					DeleteFile(mock.Anything, fileID).
					Return(nil)
				return m
			},
		},
		{
			name: "file not found",
			args: args{
				fileID: fileID,
			},
			statusCode: http.StatusBadRequest,
			serviceMock: func() service.StorageService {
				m := mocks.NewStorageService(t)
				m.EXPECT().
					DeleteFile(mock.Anything, fileID).
					Return(storage.ErrFileNotFound)
				return m
			},
		},
		{
			name: "forbidden",
			args: args{
				fileID: fileID,
			},
			statusCode: http.StatusForbidden,
			serviceMock: func() service.StorageService {
				m := mocks.NewStorageService(t)
				m.EXPECT().
					DeleteFile(mock.Anything, fileID).
					Return(storage.ErrForbidden)
				return m
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			e := echo.New()
			s := tt.serviceMock()
			req := httptest.NewRequest(http.MethodDelete, "/"+tt.args.fileID, nil)
			rec := httptest.NewRecorder()
			c := e.NewContext(req, rec)
			c.SetParamNames("id")
			c.SetParamValues(tt.args.fileID)
			handler := NewHandler(s)
			middlewares.ErrorMiddleware()(handler.DeleteFileHandler())(c)
			require.Equal(t, tt.statusCode, rec.Code)
		})
	}
}
