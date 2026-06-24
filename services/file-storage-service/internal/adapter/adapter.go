package adapter

import (
	"context"
	"storage/internal/models"
)

type FileAdapter interface {
	PushUploadingFile(ctx context.Context, metadata *models.FileMetadata) error
	PushUploadedFile(ctx context.Context, metadata *models.FileMetadata) error
	PushFailedFile(ctx context.Context, metadata *models.FileMetadata, err error) error
}

type UploadingAdapter interface {
	PushCancelUploading(ctx context.Context, fileID string) error
}
