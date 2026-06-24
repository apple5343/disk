package upload

import (
	"context"
	"storage/internal/adapter"

	"github.com/redis/go-redis/v9"
)

const (
	uploadingCancelChannel = "uploading_cancel"
)

type uploadingAdapter struct {
	rdb *redis.Client
}

func NewUploadingAdapter(rdb *redis.Client) adapter.UploadingAdapter {
	return &uploadingAdapter{
		rdb: rdb,
	}
}

func (a *uploadingAdapter) PushCancelUploading(ctx context.Context, fileID string) error {
	err := a.rdb.Publish(ctx, uploadingCancelChannel, fileID).Err()
	return err
}
