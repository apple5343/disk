package redis

import (
	"context"
	"data/internal/models"
)

func (r *fileProcessingRepository) SetProcessing(ctx context.Context, file *models.FileMetadata) error {
	key := "processing:" + file.ID
	return r.rdb.Set(ctx, key, file.UserID, 0).Err()
}

func (r *fileProcessingRepository) ProcessingIsExists(ctx context.Context, file *models.FileMetadata) (bool, error) {
	key := "processing:" + file.ID
	exists, err := r.rdb.Exists(ctx, key).Result()
	if err != nil {
		return false, err
	}
	return exists == 1, nil
}

func (r *fileProcessingRepository) DeleteProcessing(ctx context.Context, file *models.FileMetadata) error {
	key := "processing:" + file.ID
	return r.rdb.Del(ctx, key).Err()
}
