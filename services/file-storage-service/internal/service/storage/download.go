package storage

import (
	"context"
	"errors"
	"io"
	filedata "storage/internal/infrastructure/file-data"
	"storage/internal/models"
	ctxutil "storage/internal/utils/ctx"

	"github.com/apple5343/errorx"
)

func (s *storageService) DownloadFile(ctx context.Context, id string) (*models.FileMetadata, io.Reader, error) {
	userID := ctxutil.UserIDFromContext(ctx)
	if userID == "" {
		return nil, nil, ErrInvalidUserID
	}

	metadata, err := s.fileDataClient.GetFileMetadata(ctx, id)
	if err != nil {
		if errors.Is(err, filedata.ErrNotFound) {
			return nil, nil, ErrFileNotFound
		}
		return nil, nil, err
	}

	if metadata.UserID != userID {
		return nil, nil, ErrForbidden
	}

	if metadata.Status != StatusUploaded {
		return nil, nil, errorx.NewError("file not uploaded", errorx.BadRequest)
	}

	reader, err := s.storageRepository.ReadFile(ctx, metadata.StoragePath)
	if err != nil {
		return nil, nil, err
	}

	return metadata, reader, nil
}
