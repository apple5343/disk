package collector

import (
	"context"
	"data/internal/models"
	folderErrors "data/internal/service/folder"
	"data/internal/service/mocks"
	ctxutil "data/internal/utils/ctx"
	"testing"

	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

func TestDeleteFolder(t *testing.T) {
	t.Parallel()

	t.Run("success", func(t *testing.T) {
		t.Parallel()
		fileService := mocks.NewFileService(t)
		folderService := mocks.NewFolderService(t)
		folder := randomFolder()

		folderService.EXPECT().
			GetFolderByID(mock.Anything, mock.Anything).
			Return(folder, nil).
			Run(func(_ context.Context, id string) {
				require.Equal(t, folder.ID, id)
			})

		fileService.EXPECT().
			GetFilesByFolderID(mock.Anything, mock.Anything).
			Return([]*models.FileMetadata{}, nil).
			Run(func(_ context.Context, id string) {
				require.Equal(t, folder.ID, id)
			})

		folderService.EXPECT().
			GetFoldersByParentID(mock.Anything, mock.Anything).
			Return([]*models.Folder{}, nil).
			Run(func(_ context.Context, id string) {
				require.Equal(t, folder.ID, id)
			})

		folderService.EXPECT().
			DeleteFolder(mock.Anything, mock.Anything).
			Return(nil).
			Run(func(_ context.Context, id string) {
				require.Equal(t, folder.ID, id)
			})

		service := NewService(fileService, folderService)
		ctx := ctxutil.ContextWithUserID(context.Background(), folder.UserID)

		err := service.DeleteFolder(ctx, folder.ID)
		require.NoError(t, err)
	})

	t.Run("folder not found", func(t *testing.T) {
		t.Parallel()
		fileService := mocks.NewFileService(t)
		folderService := mocks.NewFolderService(t)
		folder := randomFolder()

		folderService.EXPECT().
			GetFolderByID(mock.Anything, mock.Anything).
			Return(nil, folderErrors.ErrFolderNotFound).
			Run(func(_ context.Context, id string) {
				require.Equal(t, folder.ID, id)
			})

		service := NewService(fileService, folderService)
		ctx := ctxutil.ContextWithUserID(context.Background(), folder.UserID)

		err := service.DeleteFolder(ctx, folder.ID)
		require.Error(t, err)
		require.ErrorIs(t, err, folderErrors.ErrFolderNotFound)
	})

	t.Run("folder not empty", func(t *testing.T) {
		t.Parallel()
		fileService := mocks.NewFileService(t)
		folderService := mocks.NewFolderService(t)
		folder := randomFolder()

		folderService.EXPECT().
			GetFolderByID(mock.Anything, mock.Anything).
			Return(folder, nil).
			Run(func(_ context.Context, id string) {
				require.Equal(t, folder.ID, id)
			})

		fileService.EXPECT().
			GetFilesByFolderID(mock.Anything, mock.Anything).
			Return([]*models.FileMetadata{randomFile()}, nil).
			Run(func(_ context.Context, id string) {
				require.Equal(t, folder.ID, id)
			})

		folderService.EXPECT().
			GetFoldersByParentID(mock.Anything, mock.Anything).
			Return([]*models.Folder{}, nil).
			Run(func(_ context.Context, id string) {
				require.Equal(t, folder.ID, id)
			})

		service := NewService(fileService, folderService)
		ctx := ctxutil.ContextWithUserID(context.Background(), folder.UserID)

		err := service.DeleteFolder(ctx, folder.ID)
		require.Error(t, err)
		require.ErrorIs(t, err, ErrFolderNotEmpty)
	})

	t.Run("delete root folder", func(t *testing.T) {
		t.Parallel()
		fileService := mocks.NewFileService(t)
		folderService := mocks.NewFolderService(t)
		folder := randomFolder()
		folder.IsRoot = true

		folderService.EXPECT().
			GetFolderByID(mock.Anything, mock.Anything).
			Return(folder, nil).
			Run(func(_ context.Context, id string) {
				require.Equal(t, folder.ID, id)
			})

		fileService.EXPECT().
			GetFilesByFolderID(mock.Anything, mock.Anything).
			Return([]*models.FileMetadata{}, nil).
			Run(func(_ context.Context, id string) {
				require.Equal(t, folder.ID, id)
			})

		folderService.EXPECT().
			GetFoldersByParentID(mock.Anything, mock.Anything).
			Return([]*models.Folder{}, nil).
			Run(func(_ context.Context, id string) {
				require.Equal(t, folder.ID, id)
			})

		service := NewService(fileService, folderService)
		ctx := ctxutil.ContextWithUserID(context.Background(), folder.UserID)

		err := service.DeleteFolder(ctx, folder.ID)
		require.Error(t, err)
		require.ErrorIs(t, err, ErrDeleteRoot)
	})
}
