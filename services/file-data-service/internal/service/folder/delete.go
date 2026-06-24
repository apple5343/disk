package folder

import (
	"context"
	"data/internal/repository"
	"errors"
)

func (s *folderService) DeleteFolder(ctx context.Context, id string) error {
	err := s.folderRepository.DeleteFolder(ctx, id)
	if err != nil {
		if errors.Is(err, repository.ErrNotFound) {
			return ErrFolderNotFound
		}
		return err
	}
	return nil
}
