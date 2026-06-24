package httpclient

import (
	"storage/internal/models"
	"time"
)

type Error struct {
	Message string `json:"message"`
	Status  string `json:"status"`
}

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

type Folder struct {
	ID        string    `json:"id"`
	UserID    string    `json:"user_id"`
	Name      string    `json:"name"`
	ParentID  *string   `json:"parent_id"`
	FullPath  string    `json:"full_path"`
	PathDepth int       `json:"path_depth"`
	IsRoot    bool      `json:"is_root"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

func FolderToHTTP(folder *models.Folder) *Folder {
	return &Folder{
		ID:        folder.ID,
		UserID:    folder.UserID,
		Name:      folder.Name,
		ParentID:  folder.ParentID,
		FullPath:  folder.FullPath,
		PathDepth: folder.PathDepth,
		IsRoot:    folder.IsRoot,
		CreatedAt: folder.CreatedAt,
		UpdatedAt: folder.UpdatedAt,
	}
}

func FolderFromHTTP(folder *Folder) *models.Folder {
	return &models.Folder{
		ID:        folder.ID,
		UserID:    folder.UserID,
		Name:      folder.Name,
		ParentID:  folder.ParentID,
		FullPath:  folder.FullPath,
		PathDepth: folder.PathDepth,
		IsRoot:    folder.IsRoot,
		CreatedAt: folder.CreatedAt,
		UpdatedAt: folder.UpdatedAt,
	}
}
