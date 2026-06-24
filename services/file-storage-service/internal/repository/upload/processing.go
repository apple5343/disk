package upload

import (
	"context"
	"errors"
	"storage/internal/models"

	"github.com/redis/go-redis/v9"
)

func (r *uploadProcessingRepository) SetProcessing(ctx context.Context, file *models.FileMetadata) error {
	key := "uploading:" + file.UserID + ":" + file.FullPath
	return r.rdb.Set(ctx, key, "processing", 0).Err()
}

func (r *uploadProcessingRepository) ProcessingIsExists(ctx context.Context, file *models.FileMetadata) (bool, error) {
	key := "uploading:" + file.UserID + ":" + file.FullPath
	exists, err := r.rdb.Exists(ctx, key).Result()
	if err != nil {
		return false, err
	}
	return exists == 1, nil
}

func (r *uploadProcessingRepository) DeleteProcessing(ctx context.Context, file *models.FileMetadata) error {
	key := "uploading:" + file.UserID + ":" + file.FullPath
	return r.rdb.Del(ctx, key).Err()
}

func (r *uploadProcessingRepository) UserIDByFile(ctx context.Context, id string) (string, error) {
	key := "processing:" + id
	userID, err := r.rdb.Get(ctx, key).Result()
	if err != nil {
		if errors.Is(err, redis.Nil) {
			return "", nil
		}
		return "", err
	}
	return userID, nil
}
