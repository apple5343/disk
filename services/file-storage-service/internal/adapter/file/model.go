package file

import (
	"storage/internal/models"
	"time"
)

type Message struct {
	Status string   `json:"status"`
	File   Metadata `json:"file"`
	Error  string   `json:"error"`
}

func MessageToJSON(file *models.FileMetadata, status string, err string) *Message {
	return &Message{
		Status: status,
		File:   MetadataToJSON(file),
		Error:  err,
	}
}

type Metadata struct {
	ID          string            `json:"id"`
	UserID      string            `json:"user_id"`
	StoragePath string            `json:"storage_path"`
	FolderID    string            `json:"folder_id"`
	FullPath    string            `json:"full_path"`
	FileName    string            `json:"file_name"`
	Bucket      string            `json:"bucket"`
	Size        int64             `json:"size"`
	ContentType string            `json:"content_type"`
	Tags        map[string]string `json:"tags"`
	CreatedAt   time.Time         `json:"created_at"`
	UpdatedAt   time.Time         `json:"updated_at"`
}

func MetadataToJSON(file *models.FileMetadata) Metadata {
	return Metadata{
		ID:          file.ID,
		UserID:      file.UserID,
		StoragePath: file.StoragePath,
		FolderID:    file.FolderID,
		FullPath:    file.FullPath,
		FileName:    file.FileName,
		Bucket:      file.Bucket,
		Size:        file.Size,
		ContentType: file.ContentType,
		Tags:        file.Tags,
		CreatedAt:   file.CreatedAt,
		UpdatedAt:   file.UpdatedAt,
	}
}
