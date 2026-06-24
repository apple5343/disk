package file

import (
	"context"
	"data/internal/models"
	"data/internal/repository"
	"errors"
	"time"

	"github.com/apple5343/errorx"
)

const (
	UploadingStatus = "uploading"
	UploadedStatus  = "uploaded"
	FailedStatus    = "failed"

	waitProcessInterval = 1 * time.Second
	waitProcessRetry    = 60
)

var (
	ErrMissingProcess = errors.New("missing process")
)

func (s *fileService) FileUploading(ctx context.Context, metadata *models.FileMetadata) error {
	metadata.Status = UploadingStatus
	if metadata.FolderID == "" {
		rootFolder, err := s.folderService.RootFolder(ctx, metadata.UserID)
		if err != nil {
			return errorx.NewError("save uploading: "+err.Error(), errorx.Internal)
		}
		metadata.FolderID = rootFolder.ID
	}

	if err := metadata.BeforeCreate(); err != nil {
		return err
	}

	if file, err := s.fileRepository.GetFileByPath(ctx, metadata.FullPath, metadata.UserID); nil == err {
		_ = file
		file, err = s.fileRepository.UpdateFileByPath(ctx, metadata)
		if err != nil {
			return err
		}
		err = s.fileProccessing.SetProcessing(ctx, file)
		if err != nil {
			return errorx.NewError("set file processing:"+err.Error(), errorx.Internal)
		}
		return nil
	}

	file, err := s.fileRepository.SaveFile(ctx, metadata)
	if err != nil {
		if errors.Is(err, repository.ErrAlredyExists) {
			return errorx.NewError("unexpected exist for:"+metadata.ID, errorx.Internal)
		}
		return errorx.NewError("save: "+err.Error(), errorx.Internal)
	}

	err = s.fileProccessing.SetProcessing(ctx, file)
	if err != nil {
		return errorx.NewError("set file processing:"+err.Error(), errorx.Internal)
	}
	return nil
}

func (s *fileService) FileUploaded(ctx context.Context, metadata *models.FileMetadata) error {
	if err := s.waitProcess(ctx, metadata); err != nil {
		if errors.Is(err, ErrMissingProcess) {
			return errorx.NewError("missing process", errorx.Internal)
		}
		return err
	}
	metadata.Status = UploadedStatus
	if metadata.FolderID == "" {
		rootFolder, err := s.folderService.RootFolder(ctx, metadata.UserID)
		if err != nil {
			return errorx.NewError("save uploaded: "+err.Error(), errorx.Internal)
		}
		metadata.FolderID = rootFolder.ID
	}
	_, err := s.fileRepository.UpdateFileByPath(ctx, metadata)
	if err != nil {
		return errorx.NewError("update uploaded: "+err.Error(), errorx.Internal)
	}
	err = s.fileProccessing.DeleteProcessing(ctx, metadata)
	if err != nil {
		return errorx.NewError("delete file processing:"+err.Error(), errorx.Internal)
	}
	return nil
}

func (s *fileService) FileFailed(ctx context.Context, metadata *models.FileMetadata, _ error) error {
	if err := s.waitProcess(ctx, metadata); err != nil {
		if errors.Is(err, ErrMissingProcess) {
			return errorx.NewError("missing process", errorx.Internal)
		}
		return err
	}
	metadata.Status = FailedStatus
	if metadata.FolderID == "" {
		rootFolder, err := s.folderService.RootFolder(ctx, metadata.UserID)
		if err != nil {
			return errorx.NewError("save failed: "+err.Error(), errorx.Internal)
		}
		metadata.FolderID = rootFolder.ID
	}
	_, err := s.fileRepository.UpdateFileByPath(ctx, metadata)
	if err != nil {
		return errorx.NewError("update failed: "+err.Error(), errorx.Internal)
	}
	err = s.fileProccessing.DeleteProcessing(ctx, metadata)
	if err != nil {
		return errorx.NewError("delete file processing:"+err.Error(), errorx.Internal)
	}
	return nil
}

func (s *fileService) waitProcess(ctx context.Context, file *models.FileMetadata) error {
	found := false
	for range waitProcessRetry {
		exists, err := s.fileProccessing.ProcessingIsExists(ctx, file)
		if err != nil {
			return errorx.NewError("wait process: "+err.Error(), errorx.Internal)
		}
		if exists {
			found = true
			break
		}
		time.Sleep(waitProcessInterval)
	}
	if !found {
		return errorx.NewError("missing process", errorx.Internal)
	}
	return nil
}
