package storage

import (
	"storage/internal/repository"

	"github.com/minio/minio-go/v7"
)

type storageRepository struct {
	m      *minio.Client
	bucket string
	prefix string
}

func NewRepository(m *minio.Client) repository.StorageRepository {
	return &storageRepository{
		m:      m,
		bucket: "user-files",
		prefix: "files/",
	}
}
