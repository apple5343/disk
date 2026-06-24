package models

import (
	"data/internal/models"
	"time"
)

type FileMetadata struct {
	ID          string            `json:"id"`
	UserID      string            `json:"user_id"`
	StoragePath string            `json:"storage_path"`
	FolderID    string            `json:"folder_id"`
	FileName    string            `json:"file_name"`
	FullPath    string            `json:"full_path"`
	Bucket      string            `json:"bucket"`
	Size        int64             `json:"size"`
	ContentType string            `json:"content_type"`
	Tags        map[string]string `json:"tags"`
	CreatedAt   time.Time         `json:"created_at"`
	UpdatedAt   time.Time         `json:"updated_at"`
	Status      string            `json:"status"`
}

func FileMetadataToHTTP(file *models.FileMetadata) *FileMetadata {
	return &FileMetadata{
		ID:          file.ID,
		UserID:      file.UserID,
		StoragePath: file.StoragePath,
		FolderID:    file.FolderID,
		FileName:    file.FileName,
		FullPath:    file.FullPath,
		Bucket:      file.Bucket,
		Size:        file.Size,
		ContentType: file.ContentType,
		Tags:        file.Tags,
		CreatedAt:   file.CreatedAt,
		UpdatedAt:   file.UpdatedAt,
		Status:      file.Status,
	}
}

func FileMetadataFromHTTP(file *FileMetadata) *models.FileMetadata {
	return &models.FileMetadata{
		ID:          file.ID,
		UserID:      file.UserID,
		StoragePath: file.StoragePath,
		FolderID:    file.FolderID,
		FileName:    file.FileName,
		FullPath:    file.FullPath,
		Bucket:      file.Bucket,
		Size:        file.Size,
		ContentType: file.ContentType,
		Tags:        file.Tags,
		CreatedAt:   file.CreatedAt,
		UpdatedAt:   file.UpdatedAt,
		Status:      file.Status,
	}
}
