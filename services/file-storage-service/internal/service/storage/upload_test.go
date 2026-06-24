package storage

import (
	"bytes"
	"context"
	"io"
	adapterMocks "storage/internal/adapter/mocks"
	clientMocks "storage/internal/infrastructure/file-data/mocks"
	"storage/internal/models"
	repoMocks "storage/internal/repository/mocks"
	ctxutil "storage/internal/utils/ctx"
	"testing"
	"time"

	"github.com/apple5343/errorx"
	"github.com/brianvoe/gofakeit"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

func randomFile() *models.FileMetadata {
	fileName := gofakeit.Word() + ".jpg"
	return &models.FileMetadata{
		FolderID:    "",
		FileName:    fileName,
		Bucket:      "test-bucket",
		Size:        gofakeit.Int64(),
		ContentType: "image/jpeg",
		Tags: map[string]string{
			"tag1": "value1",
			"tag2": "value2",
		},
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
		Status:    StatusUploaded,
	}
}

func TestUpload(t *testing.T) {
	t.Parallel()

	t.Run("success", func(t *testing.T) {
		t.Parallel()
		fileDataClient := clientMocks.NewClient(t)
		fileAdapter := adapterMocks.NewFileAdapter(t)
		uploadAdapter := adapterMocks.NewUploadingAdapter(t)
		storageRepo := repoMocks.NewStorageRepository(t)
		uploadRepo := repoMocks.NewUploadProcessingRepository(t)

		file := randomFile()
		fileData := []byte{1, 2, 3}
		file.Size = int64(len(fileData))
		userID := gofakeit.UUID()
		file.UserID = userID
		expectedFile := *file
		expectedFile.FullPath = "/" + file.FileName
		response := file

		uploadRepo.EXPECT().
			ProcessingIsExists(mock.Anything, mock.Anything).
			Return(false, nil).
			Run(func(_ context.Context, f *models.FileMetadata) {
				require.Equal(t, file.ID, f.ID)
			})

		uploadRepo.EXPECT().
			SetProcessing(mock.Anything, mock.Anything).
			Return(nil).
			Run(func(_ context.Context, f *models.FileMetadata) {
				require.Equal(t, file.ID, f.ID)
			})

		fileAdapter.EXPECT().
			PushUploadingFile(mock.Anything, mock.Anything).
			Return(nil).
			Run(func(_ context.Context, f *models.FileMetadata) {
				require.Equal(t, file.ID, f.ID)
			})

		storageRepo.EXPECT().
			UploadFile(mock.Anything, mock.Anything, mock.Anything).
			Return(response, nil).
			Run(func(_ context.Context, f *models.FileMetadata, data io.Reader) {
				require.Equal(t, file.ID, f.ID)
				got, err := io.ReadAll(data)
				require.NoError(t, err)
				require.Equal(t, fileData, got)
				response = f
			})

		fileAdapter.EXPECT().
			PushUploadedFile(mock.Anything, mock.Anything).
			Return(nil).
			Run(func(_ context.Context, f *models.FileMetadata) {
				require.Equal(t, file.ID, f.ID)
			})

		uploadRepo.EXPECT().
			DeleteProcessing(mock.Anything, mock.Anything).
			Return(nil).
			Run(func(_ context.Context, f *models.FileMetadata) {
				require.Equal(t, file.ID, f.ID)
			})

		service := NewService(storageRepo, fileDataClient, fileAdapter, uploadAdapter, uploadRepo)
		ctx := ctxutil.ContextWithUserID(context.Background(), userID)
		got, err := service.UploadFile(ctx, file, bytes.NewBuffer(fileData))
		require.NoError(t, err)
		require.Equal(t, expectedFile.UserID, got.UserID)
		require.Equal(t, expectedFile.FolderID, got.FolderID)
		require.Equal(t, expectedFile.FullPath, got.FullPath)
		time.Sleep((UploadingSimulationTime + 1) * time.Second)
	})
}

func TestCancelUpload(t *testing.T) {
	t.Parallel()

	t.Run("success", func(t *testing.T) {
		t.Parallel()

		userID := gofakeit.UUID()
		uploadID := gofakeit.UUID()

		uploadRepo := repoMocks.NewUploadProcessingRepository(t)
		uploadAdapter := adapterMocks.NewUploadingAdapter(t)

		uploadRepo.EXPECT().
			UserIDByFile(mock.Anything, mock.Anything).
			Return(userID, nil).
			Run(func(_ context.Context, id string) {
				require.Equal(t, uploadID, id)
			})
		uploadAdapter.EXPECT().
			PushCancelUploading(mock.Anything, mock.Anything).
			Return(nil).
			Run(func(_ context.Context, id string) {
				require.Equal(t, uploadID, id)
			})

		service := NewService(nil, nil, nil, uploadAdapter, uploadRepo)
		ctx := ctxutil.ContextWithUserID(context.Background(), userID)
		err := service.CancelUpload(ctx, uploadID)
		require.NoError(t, err)
	})

	t.Run("upload not found", func(t *testing.T) {
		t.Parallel()

		userID := gofakeit.UUID()
		uploadID := gofakeit.UUID()

		uploadRepo := repoMocks.NewUploadProcessingRepository(t)
		uploadAdapter := adapterMocks.NewUploadingAdapter(t)

		uploadRepo.EXPECT().
			UserIDByFile(mock.Anything, mock.Anything).
			Return("", nil).
			Run(func(_ context.Context, id string) {
				require.Equal(t, uploadID, id)
			})

		service := NewService(nil, nil, nil, uploadAdapter, uploadRepo)
		ctx := ctxutil.ContextWithUserID(context.Background(), userID)
		err := service.CancelUpload(ctx, uploadID)
		require.Error(t, err)
		cmnErr, ok := errorx.ToCommonError(err)
		require.True(t, ok)
		require.Equal(t, errorx.BadRequest, cmnErr.Code())
	})
}

func TestHandleCancelUpload(t *testing.T) {
	t.Parallel()

	t.Run("success", func(t *testing.T) {
		t.Parallel()
		fileDataClient := clientMocks.NewClient(t)
		fileAdapter := adapterMocks.NewFileAdapter(t)
		uploadAdapter := adapterMocks.NewUploadingAdapter(t)
		storageRepo := repoMocks.NewStorageRepository(t)
		uploadRepo := repoMocks.NewUploadProcessingRepository(t)

		file := randomFile()
		fileData := []byte{1, 2, 3}
		file.Size = int64(len(fileData))
		userID := gofakeit.UUID()
		file.UserID = userID
		expectedFile := *file
		expectedFile.FullPath = "/" + file.FileName

		uploadRepo.EXPECT().
			ProcessingIsExists(mock.Anything, mock.Anything).
			Return(false, nil).
			Run(func(_ context.Context, f *models.FileMetadata) {
				require.Equal(t, file.ID, f.ID)
			})

		uploadRepo.EXPECT().
			SetProcessing(mock.Anything, mock.Anything).
			Return(nil).
			Run(func(_ context.Context, f *models.FileMetadata) {
				require.Equal(t, file.ID, f.ID)
			})

		fileAdapter.EXPECT().
			PushUploadingFile(mock.Anything, mock.Anything).
			Return(nil).
			Run(func(_ context.Context, f *models.FileMetadata) {
				require.Equal(t, file.ID, f.ID)
			})

		storageRepo.EXPECT().
			UploadFile(mock.Anything, mock.Anything, mock.Anything).
			Return(nil, context.Canceled).
			Run(func(_ context.Context, f *models.FileMetadata, data io.Reader) {
				require.Equal(t, file.ID, f.ID)
				got, err := io.ReadAll(data)
				require.NoError(t, err)
				require.Equal(t, fileData, got)
			})

		fileAdapter.EXPECT().
			PushFailedFile(mock.Anything, mock.Anything, mock.Anything).
			Return(nil).
			Run(func(_ context.Context, f *models.FileMetadata, _ error) {
				require.Equal(t, file.ID, f.ID)
			})

		uploadRepo.EXPECT().
			DeleteProcessing(mock.Anything, mock.Anything).
			Return(nil).
			Run(func(_ context.Context, f *models.FileMetadata) {
				require.Equal(t, file.ID, f.ID)
			})

		service := NewService(storageRepo, fileDataClient, fileAdapter, uploadAdapter, uploadRepo)
		ctx := ctxutil.ContextWithUserID(context.Background(), userID)
		got, err := service.UploadFile(ctx, file, bytes.NewBuffer(fileData))
		require.NoError(t, err)
		require.Equal(t, expectedFile.UserID, got.UserID)
		require.Equal(t, expectedFile.FolderID, got.FolderID)
		require.Equal(t, expectedFile.FullPath, got.FullPath)
		time.Sleep(2 * time.Second)
		err = service.HandleCancelUpload(ctx, got.ID)
		require.NoError(t, err)

		time.Sleep(2 * time.Second)
	})
}
