package sqlx

import (
	"context"
	"data/internal/models"
	"data/internal/repository"
	repoModels "data/internal/repository/folder/models"
	sqlutil "data/internal/utils/sql"

	"github.com/google/uuid"
)

func (r *folderRepository) SaveFolder(ctx context.Context, folder *models.Folder) (*models.Folder, error) {
	repoFolder := repoModels.FolderToRepo(folder)
	repoFolder.ID = uuid.New().String()
	_, err := r.db.NamedExecContext(ctx, `INSERT INTO folders
		(id, user_id, name, parent_id, full_path, path_depth, is_root, created_at, updated_at)
		VALUES (:id, :user_id, :name, :parent_id, :full_path, :path_depth, :is_root, :created_at, :updated_at)`, repoFolder)
	if err != nil {
		if sqlutil.IsUniqueViolationSQL(err) {
			return nil, repository.ErrAlredyExists
		}
		return nil, err
	}
	return r.GetFolderByID(ctx, repoFolder.ID)
}
