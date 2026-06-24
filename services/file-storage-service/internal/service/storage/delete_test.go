package storage

import (
	"context"
	filedata "storage/internal/infrastructure/file-data"
	client "storage/internal/infrastructure/file-data/mocks"
	repo "storage/internal/repository/mocks"
	ctxutil "storage/internal/utils/ctx"
	"testing"

	"github.com/brianvoe/gofakeit"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

func TestDelete(t *testing.T) {
	t.Parallel()

	t.Run("success", func(t *testing.T) {
		t.Parallel()

		file := randomFile()
		file.ID = gofakeit.UUID()
		userID := gofakeit.UUID()
		file.UserID = userID
		ctx := ctxutil.ContextWithUserID(context.Background(), userID)

		fileID := gofakeit.UUID()
		fd := client.NewClient(t)
		repo := repo.NewStorageRepository(t)

		fd.EXPECT().
			GetFileMetadata(mock.Anything, fileID).
			Return(file, nil).
			Run(func(_ context.Context, id string) {
				require.Equal(t, fileID, id)
			})

		fd.EXPECT().
			DeleteFileMetadata(mock.Anything, fileID).
			Return(nil).
			Run(func(_ context.Context, id string) {
				require.Equal(t, fileID, id)
			})

		repo.EXPECT().
			DeleteFile(mock.Anything, file.StoragePath).
			Return(nil).
			Run(func(_ context.Context, path string) {
				require.Equal(t, file.FullPath, path)
			})

		service := NewService(repo, fd, nil, nil, nil)
		err := service.DeleteFile(ctx, fileID)
		require.NoError(t, err)
	})

	t.Run("file not found", func(t *testing.T) {
		t.Parallel()

		file := randomFile()
		file.ID = gofakeit.UUID()
		userID := gofakeit.UUID()
		file.UserID = userID
		ctx := ctxutil.ContextWithUserID(context.Background(), userID)

		fileID := gofakeit.UUID()
		fd := client.NewClient(t)

		fd.EXPECT().
			GetFileMetadata(mock.Anything, fileID).
			Return(nil, filedata.ErrNotFound).
			Run(func(_ context.Context, id string) {
				require.Equal(t, fileID, id)
			})

		service := NewService(nil, fd, nil, nil, nil)
		err := service.DeleteFile(ctx, fileID)
		require.ErrorIs(t, err, ErrFileNotFound)
	})

	t.Run("invalid user_id", func(t *testing.T) {
		service := NewService(nil, nil, nil, nil, nil)
		err := service.DeleteFile(context.Background(), gofakeit.UUID())
		require.ErrorIs(t, err, ErrInvalidUserID)
	})

	t.Run("forbidden", func(t *testing.T) {
		t.Parallel()

		file := randomFile()
		file.ID = gofakeit.UUID()
		userID := gofakeit.UUID()
		file.UserID = gofakeit.UUID()
		ctx := ctxutil.ContextWithUserID(context.Background(), userID)

		fd := client.NewClient(t)
		fd.EXPECT().
			GetFileMetadata(mock.Anything, file.ID).
			Return(file, nil).
			Run(func(_ context.Context, id string) {
				require.Equal(t, file.ID, id)
			})

		service := NewService(nil, fd, nil, nil, nil)
		err := service.DeleteFile(ctx, file.ID)
		require.ErrorIs(t, err, ErrForbidden)
	})
}
