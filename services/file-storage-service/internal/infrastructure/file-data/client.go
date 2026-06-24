package filedata

import (
	"context"
	"errors"
	"storage/internal/models"
)

var (
	ErrNotFound = errors.New("not found")
)

type Client interface {
	GetFolderByID(ctx context.Context, id string) (*models.Folder, error)
	SaveFileMetadata(ctx context.Context, metadata *models.FileMetadata) (*models.FileMetadata, error)
	GetFileMetadata(ctx context.Context, id string) (*models.FileMetadata, error)
	DeleteFileMetadata(ctx context.Context, id string) error
}
