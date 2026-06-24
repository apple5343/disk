package sqlx

import (
	"context"
	"data/internal/models"
	"data/internal/repository"
	"database/sql"
	"errors"

	repoModels "data/internal/repository/folder/models"

	"github.com/google/uuid"
)

func (r *folderRepository) GetFolderByID(ctx context.Context, id string) (*models.Folder, error) {
	var repoFolder repoModels.Folder
	_, err := uuid.Parse(id)
	if err != nil {
		return nil, repository.ErrNotFound
	}
	err = r.db.GetContext(
		ctx,
		&repoFolder,
		"SELECT id, user_id, name, parent_id, full_path, path_depth, is_root, created_at, updated_at FROM folders WHERE id = $1",
		id,
	)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, repository.ErrNotFound
		}
		return nil, err
	}
	return repoModels.FolderFromRepo(&repoFolder), nil
}

func (r *folderRepository) GetRootFolder(ctx context.Context, userID string) (*models.Folder, error) {
	var repoFolder repoModels.Folder

	err := r.db.GetContext(
		ctx,
		&repoFolder,
		"SELECT id, user_id, name, parent_id, full_path, path_depth, is_root, created_at, updated_at FROM folders WHERE user_id = $1 AND is_root = $2",
		userID,
		true,
	)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, repository.ErrNotFound
		}
		return nil, err
	}
	return repoModels.FolderFromRepo(&repoFolder), nil
}

func (r *folderRepository) GetFolderByParentID(ctx context.Context, parentID string) ([]*models.Folder, error) {
	var repoFolders []repoModels.Folder

	err := r.db.SelectContext(
		ctx,
		&repoFolders,
		"SELECT id, user_id, name, parent_id, full_path, path_depth, is_root, created_at, updated_at FROM folders WHERE parent_id = $1",
		parentID,
	)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return []*models.Folder{}, nil
		}
		return nil, err
	}
	folders := make([]*models.Folder, len(repoFolders))
	for i, repoFolder := range repoFolders {
		folders[i] = repoModels.FolderFromRepo(&repoFolder)
	}
	return folders, nil
}
