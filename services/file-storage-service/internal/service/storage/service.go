package storage

import (
	"context"
	"storage/internal/adapter"
	filedata "storage/internal/infrastructure/file-data"
	"storage/internal/repository"
	"storage/internal/service"
	"sync"

	"github.com/apple5343/errorx"
)

var (
	ErrInvalidUserID  = errorx.NewError("invalid user id", errorx.Unauthorized)
	ErrForbidden      = errorx.NewError("permission denied", errorx.Forbidden)
	ErrFileNotFound   = errorx.NewError("file not found", errorx.BadRequest)
	ErrFolderNotFound = errorx.NewError("folder not found", errorx.BadRequest)
)

type storageService struct {
	fileDataClient    filedata.Client
	fileAdapter       adapter.FileAdapter
	uploadAdapter     adapter.UploadingAdapter
	storageRepository repository.StorageRepository
	uploadRepository  repository.UploadProcessingRepository
	wg                sync.WaitGroup
	filesInProgress   *filesInProgress
}

func NewService(
	repo repository.StorageRepository,
	fileDataClient filedata.Client,
	fileAdapter adapter.FileAdapter,
	uploadAdapter adapter.UploadingAdapter,
	uploadRepo repository.UploadProcessingRepository,
) service.StorageService {
	return &storageService{
		fileDataClient:    fileDataClient,
		storageRepository: repo,
		fileAdapter:       fileAdapter,
		uploadAdapter:     uploadAdapter,
		uploadRepository:  uploadRepo,
		wg:                sync.WaitGroup{},
		filesInProgress:   newFilesInProgress(),
	}
}

func (s *storageService) Close(_ context.Context) error {
	s.wg.Wait()
	return nil
}
