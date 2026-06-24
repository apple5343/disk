package models

import (
	"data/internal/models"
	"encoding/json"
	"time"
)

type FileMetadata struct {
	ID          string    `json:"id"           db:"id"`
	UserID      string    `json:"user_id"      db:"user_id"`
	StoragePath string    `json:"storage_path" db:"storage_path"`
	FolderID    string    `json:"folder_id"    db:"folder_id"`
	FileName    string    `json:"file_name"    db:"file_name"`
	FullPath    string    `json:"full_path"    db:"full_path"`
	Bucket      string    `json:"bucket"       db:"bucket"`
	Size        int64     `json:"size"         db:"size"`
	ContentType string    `json:"content_type" db:"content_type"`
	Tags        []byte    `json:"tags"         db:"tags"`
	CreatedAt   time.Time `json:"created_at"   db:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"   db:"updated_at"`
	Status      string    `json:"status"       db:"status"`
}

func FileMetadataToRepo(file *models.FileMetadata) *FileMetadata {
	tags, _ := json.Marshal(file.Tags)

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
		CreatedAt:   file.CreatedAt,
		Tags:        tags,
		UpdatedAt:   file.UpdatedAt,
		Status:      file.Status,
	}
}

func FileMetadataFromRepo(file *FileMetadata) (*models.FileMetadata, error) {
	tags := make(map[string]string)
	err := json.Unmarshal(file.Tags, &tags)
	if err != nil {
		return nil, err
	}

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
		Tags:        tags,
		CreatedAt:   file.CreatedAt,
		UpdatedAt:   file.UpdatedAt,
		Status:      file.Status,
	}, nil
}
