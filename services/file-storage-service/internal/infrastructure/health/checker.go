package health

import (
	"context"
	"errors"
	"storage/internal/infrastructure/kafka"

	"github.com/minio/minio-go/v7"
	"github.com/redis/go-redis/v9"
)

const bucket = "user-files"

var ErrBucketNotExists = errors.New("bucket not exists")

type Checker struct {
	db     *minio.Client
	rdb    *redis.Client
	writer kafka.Producer
}

func NewChecker(db *minio.Client, rdb *redis.Client, writer kafka.Producer) *Checker {
	return &Checker{db: db, rdb: rdb, writer: writer}
}

func (c *Checker) Check(ctx context.Context) error {
	exists, err := c.db.BucketExists(ctx, bucket)
	if err != nil {
		return err
	}

	if !exists {
		return ErrBucketNotExists
	}

	err = c.rdb.Ping(ctx).Err()
	if err != nil {
		return err
	}

	err = c.writer.Ping(ctx)
	if err != nil {
		return err
	}

	return nil
}
