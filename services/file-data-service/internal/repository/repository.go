package repository

import (
	"context"
	"data/internal/models"
	"errors"
)

var (
	ErrAlredyExists = errors.New("alredy exists")
	ErrNotFound     = errors.New("not found")
)

type FileRepository interface {
	SaveFile(ctx context.Context, metadata *models.FileMetadata) (*models.FileMetadata, error)
	DeleteFile(ctx context.Context, id string) error
	UpdateFileByPath(ctx context.Context, metadata *models.FileMetadata) (*models.FileMetadata, error)
	GetFileByID(ctx context.Context, id string) (*models.FileMetadata, error)
	GetFileByStoragePath(ctx context.Context, storagePath string) (*models.FileMetadata, error)
	GetFileByPath(ctx context.Context, path string, userID string) (*models.FileMetadata, error)
	GetFilesByFolderID(ctx context.Context, folderID string) ([]*models.FileMetadata, error)
	SearchFiles(ctx context.Context, userID string, req *models.SearchRequest) ([]*models.FileMetadata, error)
}

type FileProcessingRepository interface {
	SetProcessing(ctx context.Context, file *models.FileMetadata) error
	ProcessingIsExists(ctx context.Context, file *models.FileMetadata) (bool, error)
	DeleteProcessing(ctx context.Context, file *models.FileMetadata) error
}

type FolderRepository interface {
	SaveFolder(ctx context.Context, folder *models.Folder) (*models.Folder, error)
	GetFolderByID(ctx context.Context, id string) (*models.Folder, error)
	GetFolderByParentID(ctx context.Context, parentID string) ([]*models.Folder, error)
	GetRootFolder(ctx context.Context, userID string) (*models.Folder, error)
	DeleteFolder(ctx context.Context, id string) error
}
