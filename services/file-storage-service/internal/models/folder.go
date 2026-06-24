package models

import "time"

type Folder struct {
	ID        string
	UserID    string
	Name      string
	ParentID  *string
	FullPath  string
	PathDepth int
	IsRoot    bool
	CreatedAt time.Time
	UpdatedAt time.Time
}
