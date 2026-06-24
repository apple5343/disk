package collector

import (
	"context"
	"data/internal/models"
)

func (s *collectorService) GetFolderTree(ctx context.Context, parentID string) (*models.FolderTree, error) {
	parent, err := s.folderService.GetFolderByID(ctx, parentID)
	if err != nil {
		return nil, err
	}

	folders, err := s.folderService.GetFoldersByParentID(ctx, parentID)
	if err != nil {
		return nil, err
	}

	files, err := s.fileService.GetFilesByFolderID(ctx, parentID)
	if err != nil {
		return nil, err
	}

	return &models.FolderTree{
		Folder: parent,
		Childs: folders,
		Files:  files,
	}, nil
}
