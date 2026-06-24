package upload

import (
	"storage/internal/repository"

	"github.com/redis/go-redis/v9"
)

type uploadProcessingRepository struct {
	rdb *redis.Client
}

func NewRepository(rdb *redis.Client) repository.UploadProcessingRepository {
	return &uploadProcessingRepository{
		rdb: rdb,
	}
}
