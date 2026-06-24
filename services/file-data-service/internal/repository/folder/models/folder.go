package models

import (
	"data/internal/models"
	"database/sql"
	"time"
)

type Folder struct {
	ID        string         `db:"id"         json:"id"`
	UserID    string         `db:"user_id"    json:"user_id"`
	Name      string         `db:"name"       json:"name"`
	ParentID  sql.NullString `db:"parent_id"  json:"parent_id"`
	FullPath  string         `db:"full_path"  json:"full_path"`
	PathDepth int            `db:"path_depth" json:"path_depth"`
	IsRoot    bool           `db:"is_root"    json:"is_root"`
	CreatedAt time.Time      `db:"created_at" json:"created_at"`
	UpdatedAt time.Time      `db:"updated_at" json:"updated_at"`
}

func FolderToRepo(folder *models.Folder) *Folder {
	parentID := sql.NullString{}
	if folder.ParentID != nil {
		parentID = sql.NullString{
			String: *folder.ParentID,
			Valid:  true,
		}
	}
	return &Folder{
		ID:        folder.ID,
		UserID:    folder.UserID,
		Name:      folder.Name,
		ParentID:  parentID,
		FullPath:  folder.FullPath,
		PathDepth: folder.PathDepth,
		IsRoot:    folder.IsRoot,
		CreatedAt: folder.CreatedAt,
		UpdatedAt: folder.UpdatedAt,
	}
}

func FolderFromRepo(folder *Folder) *models.Folder {
	var parentID *string
	if folder.ParentID.Valid {
		parentID = &folder.ParentID.String
	}
	return &models.Folder{
		ID:        folder.ID,
		UserID:    folder.UserID,
		Name:      folder.Name,
		ParentID:  parentID,
		FullPath:  folder.FullPath,
		PathDepth: folder.PathDepth,
		IsRoot:    folder.IsRoot,
		CreatedAt: folder.CreatedAt,
		UpdatedAt: folder.UpdatedAt,
	}
}
