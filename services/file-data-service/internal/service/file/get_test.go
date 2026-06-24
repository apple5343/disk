package file

import (
	"context"
	"data/internal/models"
	"data/internal/repository"
	repoMocks "data/internal/repository/mocks"
	folderService "data/internal/service/folder"
	"data/internal/service/mocks"
	ctxutil "data/internal/utils/ctx"
	"testing"

	"github.com/brianvoe/gofakeit"
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

func TestGetByID(t *testing.T) {
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

		service := NewService(fr, fs, nil)
		ctx := ctxutil.ContextWithUserID(context.Background(), file.UserID)

		res, err := service.GetFileByID(ctx, file.ID)
		require.NoError(t, err)
		require.Equal(t, *file, *res)
	})

	t.Run("empty user id", func(t *testing.T) {
		t.Parallel()
		service := NewService(nil, nil, nil)
		_, err := service.GetFileByID(context.Background(), gofakeit.UUID())
		require.ErrorIs(t, err, ErrInvalidToken)
	})

	t.Run("invalid file id", func(t *testing.T) {
		t.Parallel()
		service := NewService(nil, nil, nil)
		_, err := service.GetFileByID(ctxutil.ContextWithUserID(context.Background(), gofakeit.UUID()), "invalid-id")
		require.ErrorIs(t, err, ErrInvalidID)
	})
}

func TestGetByPath(t *testing.T) {
	t.Parallel()

	t.Run("success", func(t *testing.T) {
		t.Parallel()

		file := randomFile()
		fs := mocks.NewFolderService(t)
		fr := repoMocks.NewFileRepository(t)

		fr.EXPECT().
			GetFileByStoragePath(mock.Anything, mock.Anything).
			Return(file, nil).
			Run(func(_ context.Context, path string) {
				require.Equal(t, file.StoragePath, path)
			})

		service := NewService(fr, fs, nil)
		ctx := ctxutil.ContextWithUserID(context.Background(), file.UserID)

		res, err := service.GetFileByStoragePath(ctx, file.StoragePath)
		require.NoError(t, err)
		require.Equal(t, *file, *res)
	})

	t.Run("file not found", func(t *testing.T) {
		t.Parallel()

		file := randomFile()
		fs := mocks.NewFolderService(t)
		fr := repoMocks.NewFileRepository(t)

		fr.EXPECT().
			GetFileByStoragePath(mock.Anything, mock.Anything).
			Return(nil, repository.ErrNotFound).
			Run(func(_ context.Context, path string) {
				require.Equal(t, file.StoragePath, path)
			})

		service := NewService(fr, fs, nil)
		ctx := ctxutil.ContextWithUserID(context.Background(), file.UserID)

		_, err := service.GetFileByStoragePath(ctx, file.StoragePath)
		require.ErrorIs(t, err, ErrFileNotFound)
	})
}

func TestGetByFolderID(t *testing.T) {
	t.Parallel()

	t.Run("success", func(t *testing.T) {
		t.Parallel()

		folder := randomFolder()
		files := []*models.FileMetadata{randomFile(), randomFile()}
		fs := mocks.NewFolderService(t)
		fr := repoMocks.NewFileRepository(t)

		fs.EXPECT().
			GetFolderByID(mock.Anything, mock.Anything).
			Return(folder, nil).
			Run(func(_ context.Context, id string) {
				require.Equal(t, folder.ID, id)
			})

		fr.EXPECT().
			GetFilesByFolderID(mock.Anything, mock.Anything).
			Return(files, nil).
			Run(func(_ context.Context, id string) {
				require.Equal(t, folder.ID, id)
			})

		service := NewService(fr, fs, nil)
		ctx := ctxutil.ContextWithUserID(context.Background(), folder.UserID)

		res, err := service.GetFilesByFolderID(ctx, folder.ID)
		require.NoError(t, err)
		require.Len(t, res, len(files))
		for i, file := range files {
			require.Equal(t, *file, *res[i])
		}
	})

	t.Run("invalid file id", func(t *testing.T) {
		t.Parallel()
		service := NewService(nil, nil, nil)
		_, err := service.GetFilesByFolderID(
			ctxutil.ContextWithUserID(context.Background(), gofakeit.UUID()),
			"invalid-id",
		)
		require.ErrorIs(t, err, ErrInvalidID)
	})

	t.Run("empty user id", func(t *testing.T) {
		t.Parallel()
		service := NewService(nil, nil, nil)
		_, err := service.GetFilesByFolderID(context.Background(), gofakeit.UUID())
		require.ErrorIs(t, err, ErrInvalidToken)
	})

	t.Run("user id does not match", func(t *testing.T) {
		t.Parallel()

		folder := randomFolder()
		fs := mocks.NewFolderService(t)
		fr := repoMocks.NewFileRepository(t)

		fs.EXPECT().
			GetFolderByID(mock.Anything, mock.Anything).
			Return(nil, folderService.ErrFolderNotFound).
			Run(func(_ context.Context, id string) {
				require.Equal(t, folder.ID, id)
			})

		service := NewService(fr, fs, nil)
		ctx := ctxutil.ContextWithUserID(context.Background(), gofakeit.UUID())

		_, err := service.GetFilesByFolderID(ctx, folder.ID)
		require.ErrorIs(t, err, folderService.ErrFolderNotFound)
	})
}

func TestSearchFiles(t *testing.T) {
	t.Parallel()

	t.Run("success", func(t *testing.T) {
		t.Parallel()

		files := []*models.FileMetadata{randomFile(), randomFile()}
		req := &models.SearchRequest{
			FolderID: gofakeit.UUID(),
			Path:     files[0].FullPath,
			FileName: files[0].FileName,
		}
		userID := gofakeit.UUID()
		fs := mocks.NewFolderService(t)
		fr := repoMocks.NewFileRepository(t)

		fr.EXPECT().
			SearchFiles(mock.Anything, mock.Anything, mock.Anything).
			Return(files, nil).
			Run(func(_ context.Context, id string, r *models.SearchRequest) {
				require.Equal(t, userID, id)
				require.Equal(t, *req, *r)
			})

		service := NewService(fr, fs, nil)
		ctx := ctxutil.ContextWithUserID(context.Background(), userID)

		res, err := service.SearchFiles(ctx, req)
		require.NoError(t, err)
		require.Len(t, res, len(files))
		for i, file := range files {
			require.Equal(t, *file, *res[i])
		}
	})

	t.Run("invalid folder id", func(t *testing.T) {
		t.Parallel()
		service := NewService(nil, nil, nil)
		_, err := service.SearchFiles(
			ctxutil.ContextWithUserID(context.Background(), gofakeit.UUID()),
			&models.SearchRequest{FolderID: "invalid-id"},
		)
		require.ErrorIs(t, err, ErrInvalidID)
	})

	t.Run("empty user id", func(t *testing.T) {
		t.Parallel()
		service := NewService(nil, nil, nil)
		_, err := service.SearchFiles(context.Background(), &models.SearchRequest{})
		require.ErrorIs(t, err, ErrInvalidToken)
	})
}
