package storage

import (
	"context"
	"errors"
	"io"
	filedata "storage/internal/infrastructure/file-data"
	"storage/internal/models"
	"storage/internal/repository"
	ctxutil "storage/internal/utils/ctx"
	"storage/pkg/logger"
	"time"

	"github.com/apple5343/errorx"
	"github.com/google/uuid"
	"go.uber.org/zap"
)

const (
	UploadingSimulationTime = 30

	StatusUploading = "uploading"
	StatusUploaded  = "uploaded"
	StatusFailed    = "failed"
)

func (s *storageService) BeforeUploadFile(
	ctx context.Context,
	metadata *models.FileMetadata,
) (*models.FileMetadata, error) {
	userID := ctxutil.UserIDFromContext(ctx)
	if userID == "" {
		return nil, ErrInvalidUserID
	}

	metadata.BeforeCreate()
	metadata.UserID = userID

	folder := &models.Folder{
		UserID:   userID,
		FullPath: "/",
		IsRoot:   true,
	}
	var err error

	if metadata.FolderID != "" {
		folder, err = s.fileDataClient.GetFolderByID(ctx, metadata.FolderID)
	}

	if err != nil {
		if errors.Is(err, filedata.ErrNotFound) {
			return nil, ErrFolderNotFound
		}
		return nil, err
	}

	if folder.UserID != userID {
		return nil, ErrFolderNotFound
	}

	if folder.IsRoot {
		metadata.FullPath = folder.FullPath + metadata.FileName
	} else {
		metadata.FullPath = folder.FullPath + "/" + metadata.FileName
	}
	metadata.ID = uuid.New().String()

	exists, err := s.uploadRepository.ProcessingIsExists(ctx, metadata)
	if err != nil {
		return nil, errorx.NewError("upload file: "+err.Error(), errorx.Internal)
	}
	if exists {
		return nil, errorx.NewError("file already in progress", errorx.BadRequest)
	}

	return metadata, nil
}

func (s *storageService) UploadFile(
	ctx context.Context,
	metadata *models.FileMetadata,
	data io.Reader,
) (*models.FileMetadata, error) {
	metadata, err := s.BeforeUploadFile(ctx, metadata)
	if err != nil {
		return nil, err
	}

	err = s.uploadRepository.SetProcessing(ctx, metadata)
	if err != nil {
		return nil, errorx.NewError("upload file: set processing: "+err.Error(), errorx.Internal)
	}
	err = s.fileAdapter.PushUploadingFile(ctx, metadata)
	if err != nil {
		return nil, errorx.NewError("upload file: push: "+err.Error(), errorx.Internal)
	}
	l, lOk := logger.FromContext(ctx)
	if !lOk {
		l = logger.NewBaseLogger()
	}
	s.wg.Add(1)
	processingCtx, cancel := context.WithCancel(context.Background())
	s.filesInProgress.set(metadata.ID, &FileProgress{processingCtx, cancel})

	go func() {
		sleepWithContext(processingCtx, UploadingSimulationTime*time.Second) // Симуляция загрузки файла
		defer s.wg.Done()
		defer cancel()
		closer, ok := data.(io.Closer)
		if ok {
			defer closer.Close()
		}
		err = s.ProcessFile(processingCtx, metadata, data)
		if err != nil {
			if errors.Is(err, context.Canceled) {
				l.Info(ctx, "upload canceled", zap.String("id", metadata.ID))
			} else {
				l.Error(ctx, "process file:"+err.Error(), zap.String("id", metadata.ID))
			}
			err = s.fileAdapter.PushFailedFile(context.Background(), metadata, err)
			if err != nil {
				l.Error(ctx, "process file:"+err.Error(), zap.String("id", metadata.ID))
			}
		}
		err = s.uploadRepository.DeleteProcessing(context.Background(), metadata)
		if err != nil {
			l.Error(ctx, "delete processing:"+err.Error(), zap.String("id", metadata.ID))
		}
		s.filesInProgress.delete(metadata.ID)
	}()

	return metadata, nil
}

func (s *storageService) ProcessFile(ctx context.Context, metadata *models.FileMetadata, data io.Reader) error {
	file, err := s.storageRepository.UploadFile(ctx, metadata, data)
	if err != nil {
		if errors.Is(err, repository.ErrInvalidPath) {
			return errorx.NewError("invalid path: "+metadata.FullPath, errorx.BadRequest)
		}
		if errors.Is(err, context.Canceled) {
			return err
		}
		return errorx.NewError(err.Error(), errorx.Internal)
	}
	err = s.fileAdapter.PushUploadedFile(ctx, file)
	if err != nil {
		return errorx.NewError("process file: push: "+err.Error(), errorx.Internal)
	}
	return nil
}

func (s *storageService) CancelUpload(ctx context.Context, id string) error {
	userID := ctxutil.UserIDFromContext(ctx)
	if userID == "" {
		return ErrInvalidUserID
	}
	uploadingUserID, err := s.uploadRepository.UserIDByFile(ctx, id)
	if err != nil {
		return errorx.NewError("cancel upload: "+err.Error(), errorx.Internal)
	}
	if uploadingUserID != userID || uploadingUserID == "" {
		return errorx.NewError("upload not found", errorx.BadRequest)
	}
	err = s.uploadAdapter.PushCancelUploading(ctx, id)
	if err != nil {
		return errorx.NewError("cancel upload: "+err.Error(), errorx.Internal)
	}
	return nil
}

func (s *storageService) HandleCancelUpload(_ context.Context, id string) error {
	process, ok := s.filesInProgress.get(id)
	if !ok {
		return nil
	}
	process.cancel()
	return nil
}

func sleepWithContext(ctx context.Context, duration time.Duration) {
	timer := time.NewTimer(duration)
	defer timer.Stop()

	select {
	case <-timer.C:
		return
	case <-ctx.Done():
		return
	}
}
