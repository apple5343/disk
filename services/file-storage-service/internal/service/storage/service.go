package storage

import (
	"context"
	"log"
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

func (s *storageService) Close(ctx context.Context) error {
	waitCh := make(chan struct{})
	go func() {
		s.wg.Wait()
		close(waitCh)
	}()
	select {
	case <-waitCh:
		log.Println("storage service closed successfully")
	case <-ctx.Done():
		log.Println("storage service closed with timeout")
	}
	return nil
}
