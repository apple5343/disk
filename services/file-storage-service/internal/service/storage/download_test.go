package storage

import (
	"bytes"
	"context"
	"io"
	filedata "storage/internal/infrastructure/file-data"
	filedataMocks "storage/internal/infrastructure/file-data/mocks"
	repositoryMocks "storage/internal/repository/mocks"
	ctxutil "storage/internal/utils/ctx"
	"testing"

	"github.com/brianvoe/gofakeit"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

func TestDownload(t *testing.T) {
	t.Parallel()

	t.Run("success", func(t *testing.T) {
		t.Parallel()

		file := randomFile()
		file.ID = gofakeit.UUID()
		fileData := []byte{1, 2, 3}
		userID := gofakeit.UUID()
		file.UserID = userID
		ctx := ctxutil.ContextWithUserID(context.Background(), userID)

		repoMock := repositoryMocks.NewStorageRepository(t)
		filedataMock := filedataMocks.NewClient(t)

		repoMock.EXPECT().
			ReadFile(mock.Anything, mock.Anything).
			Return(bytes.NewBuffer(fileData), nil).
			Run(func(_ context.Context, path string) {
				require.Equal(t, file.FullPath, path)
			})
		filedataMock.EXPECT().
			GetFileMetadata(mock.Anything, mock.Anything).
			Return(file, nil).
			Run(func(_ context.Context, id string) {
				require.Equal(t, file.ID, id)
			})

		service := NewService(repoMock, filedataMock, nil, nil, nil)
		metadata, reader, err := service.DownloadFile(ctx, file.ID)
		require.NoError(t, err)
		require.Equal(t, file, metadata)
		got, err := io.ReadAll(reader)
		require.NoError(t, err)
		require.Equal(t, fileData, got)
	})

	t.Run("invalid user id", func(t *testing.T) {
		t.Parallel()

		service := NewService(nil, nil, nil, nil, nil)
		_, _, err := service.DownloadFile(context.Background(), gofakeit.UUID())
		require.ErrorIs(t, err, ErrInvalidUserID)
	})

	t.Run("invalid file", func(t *testing.T) {
		t.Parallel()

		file := randomFile()
		file.ID = gofakeit.UUID()
		userID := gofakeit.UUID()
		file.UserID = userID
		ctx := ctxutil.ContextWithUserID(context.Background(), userID)

		filedataMock := filedataMocks.NewClient(t)
		filedataMock.EXPECT().
			GetFileMetadata(mock.Anything, mock.Anything).
			Return(nil, filedata.ErrNotFound).
			Run(func(_ context.Context, id string) {
				require.Equal(t, file.ID, id)
			})

		service := NewService(nil, filedataMock, nil, nil, nil)
		_, _, err := service.DownloadFile(ctx, file.ID)
		require.ErrorIs(t, err, ErrFileNotFound)
	})

	t.Run("forbidden", func(t *testing.T) {
		t.Parallel()

		file := randomFile()
		file.ID = gofakeit.UUID()
		userID := gofakeit.UUID()
		file.UserID = userID
		ctx := ctxutil.ContextWithUserID(context.Background(), gofakeit.UUID())

		filedataMock := filedataMocks.NewClient(t)
		filedataMock.EXPECT().
			GetFileMetadata(mock.Anything, mock.Anything).
			Return(file, nil).
			Run(func(_ context.Context, id string) {
				require.Equal(t, file.ID, id)
			})

		service := NewService(nil, filedataMock, nil, nil, nil)
		_, _, err := service.DownloadFile(ctx, file.ID)
		require.ErrorIs(t, err, ErrForbidden)
	})
}
