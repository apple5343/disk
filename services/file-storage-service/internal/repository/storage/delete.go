package storage

import (
	"context"

	"github.com/minio/minio-go/v7"
)

func (r *storageRepository) DeleteFile(ctx context.Context, path string) error {
	err := r.m.RemoveObject(ctx, r.bucket, path, minio.RemoveObjectOptions{})
	if err != nil {
		return err
	}

	return nil
}
