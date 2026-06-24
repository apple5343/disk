package storage

import (
	"context"
	"io"
	"storage/internal/models"
	"storage/internal/repository"
	"time"

	"github.com/minio/minio-go/v7"
)

func (r *storageRepository) UploadFile(
	ctx context.Context,
	file *models.FileMetadata,
	data io.Reader,
) (*models.FileMetadata, error) {
	path := r.prefix + file.UserID + file.FullPath
	file.Bucket = r.bucket
	file.StoragePath = path
	opts := minio.PutObjectOptions{
		ContentType: file.ContentType,
		UserMetadata: map[string]string{
			"user-id":     file.UserID,
			"folder-id":   file.FolderID,
			"file-name":   file.FileName,
			"uploaded-at": file.CreatedAt.Format(time.RFC3339),
		},
		UserTags: file.Tags,
	}

	_, err := r.m.PutObject(ctx, file.Bucket, path, data, file.Size, opts)
	if err != nil {
		resp := minio.ToErrorResponse(err)
		if resp.Code == minio.XMinioInvalidObjectName {
			return nil, repository.ErrInvalidPath
		}
		return nil, err
	}

	return file, nil
}
