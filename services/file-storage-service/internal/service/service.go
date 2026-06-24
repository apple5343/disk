package service

import (
	"context"
	"io"
	"storage/internal/models"
)

type StorageService interface {
	UploadFile(ctx context.Context, metadata *models.FileMetadata, data io.Reader) (*models.FileMetadata, error)
	CancelUpload(ctx context.Context, id string) error
	HandleCancelUpload(ctx context.Context, id string) error
	DeleteFile(ctx context.Context, id string) error
	DownloadFile(ctx context.Context, id string) (*models.FileMetadata, io.Reader, error)
	Close(ctx context.Context) error
}
