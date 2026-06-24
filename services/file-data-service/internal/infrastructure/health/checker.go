package health

import (
	"context"
	"data/internal/infrastructure/kafka"

	"github.com/jmoiron/sqlx"
	"github.com/redis/go-redis/v9"
)

type Checker struct {
	db       *sqlx.DB
	rdb      *redis.Client
	consumer kafka.Consumer
}

func NewChecker(db *sqlx.DB, rdb *redis.Client, consumer kafka.Consumer) *Checker {
	return &Checker{db: db, rdb: rdb, consumer: consumer}
}

func (c *Checker) Check(ctx context.Context) error {
	err := c.db.PingContext(ctx)
	if err != nil {
		return err
	}
	err = c.rdb.Ping(ctx).Err()
	if err != nil {
		return err
	}
	err = c.consumer.Ping(ctx)
	if err != nil {
		return err
	}
	return nil
}
