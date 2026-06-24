package folder

import (
	"context"
	"data/internal/models"
	"data/internal/repository"
	ctxutil "data/internal/utils/ctx"
	"errors"

	"github.com/google/uuid"
)

func (s *folderService) GetFolderByID(ctx context.Context, id string) (*models.Folder, error) {
	userID := ctxutil.UserIDFromContext(ctx)
	if userID == "" {
		return nil, ErrInvalidToken
	}
	_, err := uuid.Parse(id)
	if err != nil {
		return nil, ErrInvalidID
	}
	folder, err := s.folderRepository.GetFolderByID(ctx, id)
	if err != nil {
		if errors.Is(err, repository.ErrNotFound) {
			return nil, ErrFolderNotFound
		}
		return nil, err
	}
	if folder.UserID != userID {
		return nil, ErrFolderNotFound
	}
	return folder, nil
}

func (s *folderService) GetFoldersByParentID(ctx context.Context, parentID string) ([]*models.Folder, error) {
	_, err := uuid.Parse(parentID)
	if err != nil {
		return nil, ErrInvalidID
	}
	folders, err := s.folderRepository.GetFolderByParentID(ctx, parentID)
	if err != nil {
		return nil, err
	}
	return folders, nil
}
