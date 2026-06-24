package sqlx

import (
	"context"
	"data/internal/models"
	repository "data/internal/repository"
	repoModels "data/internal/repository/file/models"
	sqlutil "data/internal/utils/sql"
)

func (r *fileRepository) SaveFile(ctx context.Context, metadata *models.FileMetadata) (*models.FileMetadata, error) {
	repoFile := repoModels.FileMetadataToRepo(metadata)
	_, err := r.db.NamedExecContext(ctx, `INSERT INTO files
		(id, user_id, storage_path, file_name, folder_id, bucket, full_path, size, content_type, tags, created_at, updated_at, status)
		VALUES (:id, :user_id, :storage_path, :file_name, :folder_id, :bucket, :full_path, :size, :content_type, :tags, :created_at, :updated_at, :status)`, repoFile)
	if err != nil {
		if sqlutil.IsUniqueViolationSQL(err) {
			return nil, repository.ErrAlredyExists
		}
		return nil, err
	}
	return r.GetFileByID(ctx, repoFile.ID)
}
