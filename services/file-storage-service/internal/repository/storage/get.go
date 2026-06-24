package storage

import (
	"context"
	"io"

	"github.com/minio/minio-go/v7"
)

func (r *storageRepository) ReadFile(ctx context.Context, path string) (io.Reader, error) {
	object, err := r.m.GetObject(ctx, r.bucket, path, minio.GetObjectOptions{})

	if err != nil {
		return nil, err
	}

	return object, nil
}
