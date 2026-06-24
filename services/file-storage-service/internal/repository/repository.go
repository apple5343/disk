package repository

import (
	"context"
	"errors"
	"io"
	"storage/internal/models"
)

var (
	ErrFileNotFound = errors.New("file not found")
	ErrAlredyExists = errors.New("already exists")
	ErrInvalidPath  = errors.New("invalid path")
)

type StorageRepository interface {
	UploadFile(ctx context.Context, metadata *models.FileMetadata, data io.Reader) (*models.FileMetadata, error)
	DeleteFile(ctx context.Context, path string) error
	ReadFile(ctx context.Context, path string) (io.Reader, error)
}

type UploadProcessingRepository interface {
	SetProcessing(ctx context.Context, file *models.FileMetadata) error
	ProcessingIsExists(ctx context.Context, file *models.FileMetadata) (bool, error)
	DeleteProcessing(ctx context.Context, file *models.FileMetadata) error
	UserIDByFile(ctx context.Context, id string) (string, error)
}
