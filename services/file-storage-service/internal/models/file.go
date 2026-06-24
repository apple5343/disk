package models

import (
	"path"
	"time"
)

type FileMetadata struct {
	ID          string
	UserID      string
	StoragePath string
	FolderID    string
	FullPath    string
	FileName    string
	Bucket      string
	Size        int64
	ContentType string
	Tags        map[string]string
	CreatedAt   time.Time
	UpdatedAt   time.Time
	Status      string
}

func (f *FileMetadata) BeforeCreate() {
	f.FileName = path.Base(f.FileName)
	f.CreatedAt = time.Now()
	f.UpdatedAt = time.Now()
}
