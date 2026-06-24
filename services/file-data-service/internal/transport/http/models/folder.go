package models

import (
	"data/internal/models"
	"time"
)

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

type FolderTree struct {
	Folders []*Folder       `json:"folders"`
	Files   []*FileMetadata `json:"files"`
}

func FolderTreeToHTTP(tree *models.FolderTree) *FolderTree {
	folders := make([]*Folder, len(tree.Childs))
	files := make([]*FileMetadata, len(tree.Files))

	for i, f := range tree.Childs {
		folders[i] = FolderToHTTP(f)
	}
	for i, f := range tree.Files {
		files[i] = FileMetadataToHTTP(f)
	}
	return &FolderTree{
		Folders: folders,
		Files:   files,
	}
}
