package collector

import (
	"context"
	ctxutil "data/internal/utils/ctx"
)

func (s *collectorService) DeleteFolder(ctx context.Context, id string) error {
	userID := ctxutil.UserIDFromContext(ctx)
	if userID == "" {
		return ErrInvalidToken
	}

	tree, err := s.GetFolderTree(ctx, id)
	if err != nil {
		return err
	}
	if tree.Folder.IsRoot {
		return ErrDeleteRoot
	}

	if (len(tree.Childs) != 0) || (len(tree.Files) != 0) {
		return ErrFolderNotEmpty
	}
	return s.folderService.DeleteFolder(ctx, id)
}
