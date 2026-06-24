package file

import (
	"context"
	"data/internal/models"
	"data/internal/repository"
	"data/internal/repository/mocks"
	serviceMocks "data/internal/service/mocks"
	ctxutil "data/internal/utils/ctx"
	"errors"
	"testing"

	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

func TestUploading(t *testing.T) {
	t.Parallel()

	t.Run("success uploading", func(t *testing.T) {
		t.Parallel()

		file := randomFile()
		fileRepo := mocks.NewFileRepository(t)
		fileProcessRepo := mocks.NewFileProcessingRepository(t)

		fileRepo.EXPECT().
			GetFileByPath(mock.Anything, mock.Anything, mock.Anything).
			Return(nil, repository.ErrNotFound).
			Run(func(_ context.Context, path string, userID string) {
				require.Equal(t, file.UserID, userID)
				require.Equal(t, file.FullPath, path)
			})
		fileRepo.EXPECT().
			SaveFile(mock.Anything, mock.Anything).
			Return(file, nil).
			Run(func(_ context.Context, f *models.FileMetadata) {
				require.Equal(t, file.ID, f.ID)
				require.Equal(t, UploadingStatus, file.Status)
			})
		fileProcessRepo.EXPECT().
			SetProcessing(mock.Anything, mock.Anything).
			Return(nil).
			Run(func(_ context.Context, f *models.FileMetadata) {
				require.Equal(t, file.ID, f.ID)
				require.Equal(t, UploadingStatus, file.Status)
				require.Equal(t, file.UserID, f.UserID)
			})

		service := NewService(fileRepo, nil, fileProcessRepo)
		ctx := ctxutil.ContextWithUserID(context.Background(), file.UserID)

		err := service.FileUploading(ctx, file)
		require.NoError(t, err)
	})

	t.Run("success uploading without folder", func(t *testing.T) {
		t.Parallel()

		file := randomFile()
		folder := randomFolder()
		file.FolderID = ""
		fileRepo := mocks.NewFileRepository(t)
		fileProcessRepo := mocks.NewFileProcessingRepository(t)
		folderService := serviceMocks.NewFolderService(t)

		folderService.EXPECT().
			RootFolder(mock.Anything, mock.Anything).
			Return(folder, nil).
			Run(func(_ context.Context, userID string) {
				require.Equal(t, file.UserID, userID)
			})

		fileRepo.EXPECT().
			GetFileByPath(mock.Anything, mock.Anything, mock.Anything).
			Return(nil, repository.ErrNotFound).
			Run(func(_ context.Context, path string, userID string) {
				require.Equal(t, file.UserID, userID)
				require.Equal(t, file.FullPath, path)
			})
		fileRepo.EXPECT().
			SaveFile(mock.Anything, mock.Anything).
			Return(file, nil).
			Run(func(_ context.Context, f *models.FileMetadata) {
				require.Equal(t, file.ID, f.ID)
				require.Equal(t, folder.ID, f.FolderID)
				require.Equal(t, UploadingStatus, file.Status)
			})
		fileProcessRepo.EXPECT().
			SetProcessing(mock.Anything, mock.Anything).
			Return(nil).
			Run(func(_ context.Context, f *models.FileMetadata) {
				require.Equal(t, file.ID, f.ID)
				require.Equal(t, UploadingStatus, file.Status)
				require.Equal(t, file.UserID, f.UserID)
			})

		service := NewService(fileRepo, folderService, fileProcessRepo)
		ctx := ctxutil.ContextWithUserID(context.Background(), file.UserID)

		err := service.FileUploading(ctx, file)
		require.NoError(t, err)
	})

	t.Run("success uploading with existing file", func(t *testing.T) {
		t.Parallel()

		file := randomFile()
		fileRepo := mocks.NewFileRepository(t)
		fileProcessRepo := mocks.NewFileProcessingRepository(t)

		fileRepo.EXPECT().
			GetFileByPath(mock.Anything, mock.Anything, mock.Anything).
			Return(file, nil).
			Run(func(_ context.Context, path string, userID string) {
				require.Equal(t, file.UserID, userID)
				require.Equal(t, file.FullPath, path)
			})
		fileRepo.EXPECT().
			UpdateFileByPath(mock.Anything, mock.Anything).
			Return(file, nil).
			Run(func(_ context.Context, f *models.FileMetadata) {
				require.Equal(t, file.ID, f.ID)
				require.Equal(t, file.UserID, f.UserID)
				require.Equal(t, file.FullPath, f.FullPath)
				require.Equal(t, UploadingStatus, file.Status)
			})
		fileProcessRepo.EXPECT().
			SetProcessing(mock.Anything, mock.Anything).
			Return(nil).
			Run(func(_ context.Context, f *models.FileMetadata) {
				require.Equal(t, file.ID, f.ID)
				require.Equal(t, UploadingStatus, file.Status)
				require.Equal(t, file.UserID, f.UserID)
			})

		service := NewService(fileRepo, nil, fileProcessRepo)
		ctx := ctxutil.ContextWithUserID(context.Background(), file.UserID)

		err := service.FileUploading(ctx, file)
		require.NoError(t, err)
	})

	t.Run("unexpected existing file", func(t *testing.T) {
		t.Parallel()

		file := randomFile()
		fileRepo := mocks.NewFileRepository(t)
		fileProcessRepo := mocks.NewFileProcessingRepository(t)

		fileRepo.EXPECT().
			GetFileByPath(mock.Anything, mock.Anything, mock.Anything).
			Return(nil, repository.ErrNotFound).
			Run(func(_ context.Context, path string, userID string) {
				require.Equal(t, file.UserID, userID)
				require.Equal(t, file.FullPath, path)
			})
		fileRepo.EXPECT().
			SaveFile(mock.Anything, mock.Anything).
			Return(nil, repository.ErrAlredyExists).
			Run(func(_ context.Context, f *models.FileMetadata) {
				require.Equal(t, file.ID, f.ID)
				require.Equal(t, UploadingStatus, file.Status)
			})
		service := NewService(fileRepo, nil, fileProcessRepo)
		ctx := ctxutil.ContextWithUserID(context.Background(), file.UserID)

		err := service.FileUploading(ctx, file)
		require.Error(t, err)
		require.Contains(t, err.Error(), "unexpected exist", file.ID)
	})
}

func TestUploaded(t *testing.T) {
	t.Parallel()

	t.Run("success", func(t *testing.T) {
		t.Parallel()

		file := randomFile()
		fileRepo := mocks.NewFileRepository(t)
		fileProcessRepo := mocks.NewFileProcessingRepository(t)

		fileProcessRepo.EXPECT().
			ProcessingIsExists(mock.Anything, file).
			Return(false, nil).
			Run(func(_ context.Context, f *models.FileMetadata) {
				require.Equal(t, file.ID, f.ID)
			}).Times(3)

		fileProcessRepo.EXPECT().
			ProcessingIsExists(mock.Anything, file).
			Return(true, nil).
			Run(func(_ context.Context, f *models.FileMetadata) {
				require.Equal(t, file.ID, f.ID)
			}).Times(1)

		fileRepo.EXPECT().
			UpdateFileByPath(mock.Anything, mock.Anything).
			Return(file, nil).
			Run(func(_ context.Context, f *models.FileMetadata) {
				require.Equal(t, file.ID, f.ID)
				require.Equal(t, file.UserID, f.UserID)
				require.Equal(t, file.FullPath, f.FullPath)
				require.Equal(t, UploadedStatus, file.Status)
			})

		fileProcessRepo.EXPECT().
			DeleteProcessing(mock.Anything, mock.Anything).
			Return(nil).
			Run(func(_ context.Context, f *models.FileMetadata) {
				require.Equal(t, file.ID, f.ID)
			})

		service := NewService(fileRepo, nil, fileProcessRepo)
		ctx := ctxutil.ContextWithUserID(context.Background(), file.UserID)

		err := service.FileUploaded(ctx, file)
		require.NoError(t, err)
	})

	t.Run("success without folder", func(t *testing.T) {
		t.Parallel()

		file := randomFile()
		file.FolderID = ""
		folder := randomFolder()
		fileRepo := mocks.NewFileRepository(t)
		fileProcessRepo := mocks.NewFileProcessingRepository(t)
		folderService := serviceMocks.NewFolderService(t)

		fileProcessRepo.EXPECT().
			ProcessingIsExists(mock.Anything, file).
			Return(false, nil).
			Run(func(_ context.Context, f *models.FileMetadata) {
				require.Equal(t, file.ID, f.ID)
			}).Times(3)

		fileProcessRepo.EXPECT().
			ProcessingIsExists(mock.Anything, file).
			Return(true, nil).
			Run(func(_ context.Context, f *models.FileMetadata) {
				require.Equal(t, file.ID, f.ID)
			}).Times(1)

		folderService.EXPECT().
			RootFolder(mock.Anything, mock.Anything).
			Return(folder, nil).
			Run(func(_ context.Context, userID string) {
				require.Equal(t, file.UserID, userID)
			})

		fileRepo.EXPECT().
			UpdateFileByPath(mock.Anything, mock.Anything).
			Return(file, nil).
			Run(func(_ context.Context, f *models.FileMetadata) {
				require.Equal(t, file.ID, f.ID)
				require.Equal(t, file.UserID, f.UserID)
				require.Equal(t, file.FullPath, f.FullPath)
				require.Equal(t, folder.ID, file.FolderID)
				require.Equal(t, UploadedStatus, file.Status)
			})

		fileProcessRepo.EXPECT().
			DeleteProcessing(mock.Anything, mock.Anything).
			Return(nil).
			Run(func(_ context.Context, f *models.FileMetadata) {
				require.Equal(t, file.ID, f.ID)
			})

		service := NewService(fileRepo, folderService, fileProcessRepo)
		ctx := ctxutil.ContextWithUserID(context.Background(), file.UserID)

		err := service.FileUploaded(ctx, file)
		require.NoError(t, err)
	})

	t.Run("missing file", func(t *testing.T) {
		t.Parallel()

		file := randomFile()
		fileRepo := mocks.NewFileRepository(t)
		fileProcessRepo := mocks.NewFileProcessingRepository(t)

		fileProcessRepo.EXPECT().
			ProcessingIsExists(mock.Anything, file).
			Return(false, nil).
			Run(func(_ context.Context, f *models.FileMetadata) {
				require.Equal(t, file.ID, f.ID)
			}).Times(waitProcessRetry)

		service := NewService(fileRepo, nil, fileProcessRepo)
		ctx := ctxutil.ContextWithUserID(context.Background(), file.UserID)

		err := service.FileUploaded(ctx, file)
		require.Error(t, err)
		require.Contains(t, err.Error(), "missing process")
	})
}

func TestFailed(t *testing.T) {
	t.Parallel()

	t.Run("success", func(t *testing.T) {
		t.Parallel()

		file := randomFile()
		fileRepo := mocks.NewFileRepository(t)
		fileProcessRepo := mocks.NewFileProcessingRepository(t)

		fileProcessRepo.EXPECT().
			ProcessingIsExists(mock.Anything, file).
			Return(false, nil).
			Run(func(_ context.Context, f *models.FileMetadata) {
				require.Equal(t, file.ID, f.ID)
			}).Times(3)

		fileProcessRepo.EXPECT().
			ProcessingIsExists(mock.Anything, file).
			Return(true, nil).
			Run(func(_ context.Context, f *models.FileMetadata) {
				require.Equal(t, file.ID, f.ID)
			}).Times(1)

		fileRepo.EXPECT().
			UpdateFileByPath(mock.Anything, mock.Anything).
			Return(file, nil).
			Run(func(_ context.Context, f *models.FileMetadata) {
				require.Equal(t, file.ID, f.ID)
				require.Equal(t, file.UserID, f.UserID)
				require.Equal(t, file.FullPath, f.FullPath)
				require.Equal(t, FailedStatus, file.Status)
			})

		fileProcessRepo.EXPECT().
			DeleteProcessing(mock.Anything, mock.Anything).
			Return(nil).
			Run(func(_ context.Context, f *models.FileMetadata) {
				require.Equal(t, file.ID, f.ID)
			})

		service := NewService(fileRepo, nil, fileProcessRepo)
		ctx := ctxutil.ContextWithUserID(context.Background(), file.UserID)

		err := service.FileFailed(ctx, file, errors.New("failed"))
		require.NoError(t, err)
	})

	t.Run("success without folder", func(t *testing.T) {
		t.Parallel()

		file := randomFile()
		file.FolderID = ""
		folder := randomFolder()
		fileRepo := mocks.NewFileRepository(t)
		fileProcessRepo := mocks.NewFileProcessingRepository(t)
		folderService := serviceMocks.NewFolderService(t)

		fileProcessRepo.EXPECT().
			ProcessingIsExists(mock.Anything, file).
			Return(false, nil).
			Run(func(_ context.Context, f *models.FileMetadata) {
				require.Equal(t, file.ID, f.ID)
			}).Times(3)

		fileProcessRepo.EXPECT().
			ProcessingIsExists(mock.Anything, file).
			Return(true, nil).
			Run(func(_ context.Context, f *models.FileMetadata) {
				require.Equal(t, file.ID, f.ID)
			}).Times(1)

		folderService.EXPECT().
			RootFolder(mock.Anything, mock.Anything).
			Return(folder, nil).
			Run(func(_ context.Context, userID string) {
				require.Equal(t, file.UserID, userID)
			})

		fileRepo.EXPECT().
			UpdateFileByPath(mock.Anything, mock.Anything).
			Return(file, nil).
			Run(func(_ context.Context, f *models.FileMetadata) {
				require.Equal(t, file.ID, f.ID)
				require.Equal(t, file.UserID, f.UserID)
				require.Equal(t, file.FullPath, f.FullPath)
				require.Equal(t, folder.ID, file.FolderID)
				require.Equal(t, FailedStatus, file.Status)
			})

		fileProcessRepo.EXPECT().
			DeleteProcessing(mock.Anything, mock.Anything).
			Return(nil).
			Run(func(_ context.Context, f *models.FileMetadata) {
				require.Equal(t, file.ID, f.ID)
			})

		service := NewService(fileRepo, folderService, fileProcessRepo)
		ctx := ctxutil.ContextWithUserID(context.Background(), file.UserID)

		err := service.FileFailed(ctx, file, errors.New("failed"))
		require.NoError(t, err)
	})

	t.Run("missing file", func(t *testing.T) {
		t.Parallel()

		file := randomFile()
		fileRepo := mocks.NewFileRepository(t)
		fileProcessRepo := mocks.NewFileProcessingRepository(t)

		fileProcessRepo.EXPECT().
			ProcessingIsExists(mock.Anything, file).
			Return(false, nil).
			Run(func(_ context.Context, f *models.FileMetadata) {
				require.Equal(t, file.ID, f.ID)
			}).Times(waitProcessRetry)

		service := NewService(fileRepo, nil, fileProcessRepo)
		ctx := ctxutil.ContextWithUserID(context.Background(), file.UserID)

		err := service.FileFailed(ctx, file, errors.New("failed"))
		require.Error(t, err)
		require.Contains(t, err.Error(), "missing process")
	})
}
