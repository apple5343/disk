package service

import (
	"context"
	"data/internal/models"
)

type FileService interface {
	FileUploading(ctx context.Context, metadata *models.FileMetadata) error
	FileUploaded(ctx context.Context, metadata *models.FileMetadata) error
	FileFailed(ctx context.Context, metadata *models.FileMetadata, err error) error

	SaveFile(ctx context.Context, metadata *models.FileMetadata) (*models.FileMetadata, error)
	DeleteFile(ctx context.Context, id string) error
	GetFileByID(ctx context.Context, id string) (*models.FileMetadata, error)
	GetFilesByFolderID(ctx context.Context, folderID string) ([]*models.FileMetadata, error)
	GetFileByStoragePath(ctx context.Context, storagePath string) (*models.FileMetadata, error)
	SearchFiles(ctx context.Context, req *models.SearchRequest) ([]*models.FileMetadata, error)
}

type FolderService interface {
	GetFoldersByParentID(ctx context.Context, parentID string) ([]*models.Folder, error)
	SaveFolder(ctx context.Context, folder *models.Folder) (*models.Folder, error)
	GetFolderByID(ctx context.Context, id string) (*models.Folder, error)
	RootFolder(ctx context.Context, userID string) (*models.Folder, error)
	DeleteFolder(ctx context.Context, id string) error
}

type CollectorService interface {
	GetFolderTree(ctx context.Context, parentID string) (*models.FolderTree, error)
	DeleteFolder(ctx context.Context, id string) error
}
