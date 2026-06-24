package redis

import (
	"data/internal/repository"

	"github.com/redis/go-redis/v9"
)

type fileProcessingRepository struct {
	rdb *redis.Client
}

func NewRepository(rdb *redis.Client) repository.FileProcessingRepository {
	return &fileProcessingRepository{
		rdb: rdb,
	}
}
