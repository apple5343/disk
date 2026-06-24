package sqlx

import (
	"context"
	"data/internal/models"
	"encoding/json"
	"time"
)

func (r *fileRepository) UpdateFileByPath(
	ctx context.Context,
	file *models.FileMetadata,
) (*models.FileMetadata, error) {
	file.UpdatedAt = time.Now()
	tags, err := json.Marshal(file.Tags)
	if err != nil {
		return nil, err
	}
	_, err = r.db.ExecContext(
		ctx,
		`UPDATE files SET 
		id = $1, folder_id = $2, file_name = $3, storage_path = $4, bucket = $5, size = $6, content_type = $7, tags = $8, updated_at = $9, status = $10
		WHERE full_path = $11 AND user_id = $12`,
		file.ID,
		file.FolderID,
		file.FileName,
		file.StoragePath,
		file.Bucket,
		file.Size,
		file.ContentType,
		tags,
		file.UpdatedAt,
		file.Status,
		file.FullPath,
		file.UserID,
	)
	if err != nil {
		return nil, err
	}

	return r.GetFileByID(ctx, file.ID)
}
