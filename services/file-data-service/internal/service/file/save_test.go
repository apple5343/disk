package file

import (
	"context"
	"data/internal/models"
	"data/internal/repository"
	repoMocks "data/internal/repository/mocks"
	"data/internal/service/mocks"
	ctxutil "data/internal/utils/ctx"
	"testing"

	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

func TestSave(t *testing.T) {
	t.Parallel()

	t.Run("success", func(t *testing.T) {
		t.Parallel()

		file := randomFile()
		file.FolderID = ""
		folder := randomFolder()
		folder.IsRoot = true
		folder.ParentID = nil
		folder.UserID = file.UserID
		folder.FullPath = "/"
		folder.PathDepth = 1
		file.FullPath = "/" + file.FileName
		fs := mocks.NewFolderService(t)
		fr := repoMocks.NewFileRepository(t)

		fs.EXPECT().
			RootFolder(mock.Anything, mock.Anything).
			Return(folder, nil).
			Run(func(_ context.Context, userID string) {
				require.Equal(t, file.UserID, userID)
			})
		fr.EXPECT().
			SaveFile(mock.Anything, mock.Anything).
			Return(file, nil).
			Run(func(_ context.Context, f *models.FileMetadata) {
				require.Equal(t, "/"+file.FileName, f.FullPath)
			})

		service := NewService(fr, fs, nil)
		ctx := ctxutil.ContextWithUserID(context.Background(), file.UserID)

		res, err := service.SaveFile(ctx, file)
		require.NoError(t, err)
		require.Equal(t, *file, *res)
	})

	t.Run("save with update", func(t *testing.T) {
		t.Parallel()

		file := randomFile()
		folder := randomFolder()
		file.FolderID = folder.ID

		fs := mocks.NewFolderService(t)
		fr := repoMocks.NewFileRepository(t)

		fr.EXPECT().
			SaveFile(mock.Anything, mock.Anything).
			Return(nil, repository.ErrAlredyExists).
			Run(func(_ context.Context, f *models.FileMetadata) {
				require.Equal(t, file.ID, f.ID)
			})

		fr.EXPECT().
			GetFileByStoragePath(mock.Anything, mock.Anything).
			Return(file, nil).
			Run(func(_ context.Context, path string) {
				require.Equal(t, file.StoragePath, path)
			})

		fr.EXPECT().
			UpdateFileByPath(mock.Anything, mock.Anything).
			Return(file, nil).
			Run(func(_ context.Context, f *models.FileMetadata) {
				require.Equal(t, file.ID, f.ID)
			})

		service := NewService(fr, fs, nil)
		ctx := ctxutil.ContextWithUserID(context.Background(), file.UserID)

		res, err := service.SaveFile(ctx, file)
		require.NoError(t, err)
		require.Equal(t, *file, *res)
	})
}
