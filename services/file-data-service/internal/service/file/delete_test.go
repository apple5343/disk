package file

import (
	"context"
	"data/internal/models"
	"data/internal/repository"
	"data/internal/service/mocks"
	ctxutil "data/internal/utils/ctx"
	"testing"
	"time"

	repoMocks "data/internal/repository/mocks"

	"github.com/brianvoe/gofakeit"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

func randomFile() *models.FileMetadata {
	folder := "/" + gofakeit.Word() + "/" + gofakeit.Word()
	fileName := gofakeit.Word() + ".jpg"
	fullPath := folder + "/" + fileName
	return &models.FileMetadata{
		ID:          gofakeit.UUID(),
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
		UpdatedAt: time.Now(),
	}
}

func TestDelete(t *testing.T) {
	t.Parallel()

	t.Run("success", func(t *testing.T) {
		t.Parallel()

		file := randomFile()
		fs := mocks.NewFolderService(t)
		fr := repoMocks.NewFileRepository(t)

		fr.EXPECT().
			GetFileByID(mock.Anything, mock.Anything).
			Return(file, nil).
			Run(func(_ context.Context, id string) {
				require.Equal(t, file.ID, id)
			})

		fr.EXPECT().
			DeleteFile(mock.Anything, mock.Anything).
			Return(nil).
			Run(func(_ context.Context, id string) {
				require.Equal(t, file.ID, id)
			})

		service := NewService(fr, fs, nil)
		ctx := ctxutil.ContextWithUserID(context.Background(), file.UserID)

		err := service.DeleteFile(ctx, file.ID)
		require.NoError(t, err)
	})

	t.Run("file not found", func(t *testing.T) {
		t.Parallel()

		fr := repoMocks.NewFileRepository(t)
		fs := mocks.NewFolderService(t)

		folder := randomFile()

		fr.EXPECT().
			GetFileByID(mock.Anything, mock.Anything).
			Return(nil, repository.ErrNotFound).
			Run(func(_ context.Context, id string) {
				require.Equal(t, folder.ID, id)
			})

		service := NewService(fr, fs, nil)
		ctx := ctxutil.ContextWithUserID(context.Background(), folder.UserID)

		err := service.DeleteFile(ctx, folder.ID)
		require.ErrorIs(t, err, ErrFileNotFound)
	})

	t.Run("empty user id", func(t *testing.T) {
		service := NewService(nil, nil, nil)
		err := service.DeleteFile(context.Background(), gofakeit.UUID())
		require.ErrorIs(t, err, ErrInvalidToken)
	})

	t.Run("invalid file id", func(t *testing.T) {
		service := NewService(nil, nil, nil)
		err := service.DeleteFile(ctxutil.ContextWithUserID(context.Background(), gofakeit.UUID()), "invalid-id")
		require.ErrorIs(t, err, ErrInvalidID)
	})

	t.Run("user id does not match", func(t *testing.T) {
		t.Parallel()

		file := randomFile()
		userID := gofakeit.UUID()
		fs := mocks.NewFolderService(t)
		fr := repoMocks.NewFileRepository(t)

		fr.EXPECT().
			GetFileByID(mock.Anything, mock.Anything).
			Return(file, nil).
			Run(func(_ context.Context, id string) {
				require.Equal(t, file.ID, id)
			})

		service := NewService(fr, fs, nil)
		ctx := ctxutil.ContextWithUserID(context.Background(), userID)

		err := service.DeleteFile(ctx, file.ID)
		require.ErrorIs(t, err, ErrFileNotFound)
	})
}
