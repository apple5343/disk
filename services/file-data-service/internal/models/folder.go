package models

import (
	"errors"
	"strings"
	"time"

	"github.com/go-playground/validator/v10"
)

var (
	ErrInvalidName = errors.New("invalid name")
)

type Folder struct {
	ID        string
	UserID    string  `validate:"required"`
	Name      string  `validate:"required,max=255"`
	ParentID  *string `validate:"required"`
	FullPath  string
	PathDepth int
	IsRoot    bool
	CreatedAt time.Time
	UpdatedAt time.Time
}

func (f *Folder) BeforeCreate() error {
	f.CreatedAt = time.Now()
	f.UpdatedAt = time.Now()
	f.Name = strings.TrimSpace(f.Name)
	f.Name = strings.Trim(f.Name, "/")
	if strings.Contains(f.Name, "/") {
		return ErrInvalidName
	}
	if (f.Name == ".") || (f.Name == "..") {
		return ErrInvalidName
	}
	if strings.ContainsAny(f.Name, `\:*?"<>|`) {
		return ErrInvalidName
	}
	err := validator.New().Struct(*f)
	if err != nil {
		return err
	}
	return nil
}

type FolderTree struct {
	Folder *Folder
	Childs []*Folder
	Files  []*FileMetadata
}
