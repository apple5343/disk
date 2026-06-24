package storage

import (
	"context"
	"errors"
	filedata "storage/internal/infrastructure/file-data"
	ctxutil "storage/internal/utils/ctx"

	"github.com/apple5343/errorx"
)

func (s *storageService) DeleteFile(ctx context.Context, id string) error {
	userID := ctxutil.UserIDFromContext(ctx)
	if userID == "" {
		return ErrInvalidUserID
	}
	file, err := s.fileDataClient.GetFileMetadata(ctx, id)
	if err != nil {
		if errors.Is(err, filedata.ErrNotFound) {
			return ErrFileNotFound
		}
		return err
	}
	if file.UserID != userID {
		return ErrForbidden
	}
	if file.Status != StatusUploaded {
		return errorx.NewError("file not loaded", errorx.BadRequest)
	}
	err = s.fileDataClient.DeleteFileMetadata(ctx, id)
	if err != nil {
		if errors.Is(err, filedata.ErrNotFound) {
			return ErrFileNotFound
		}
		return err
	}

	return s.storageRepository.DeleteFile(ctx, file.StoragePath)
}
