package models

import (
	"path"
	"time"

	"github.com/go-playground/validator/v10"
)

type FileMetadata struct {
	ID          string
	UserID      string `validate:"required"`
	StoragePath string
	FolderID    string `validate:"required"`
	FullPath    string `validate:"required"`
	FileName    string `validate:"required"`
	Bucket      string
	Size        int64
	ContentType string `validate:"required"`
	Tags        map[string]string
	CreatedAt   time.Time
	UpdatedAt   time.Time
	Status      string
}

func (f *FileMetadata) BeforeCreate() error {
	f.FileName = path.Base(f.FileName)
	f.CreatedAt = time.Now()
	f.UpdatedAt = f.CreatedAt
	err := validator.New().Struct(*f)
	if err != nil {
		return err
	}
	return nil
}
