package kafka

import (
	"context"
	"data/internal/config"
	"data/internal/infrastructure/kafka"
	"data/internal/service"
	"sync"

	"golang.org/x/sync/semaphore"
	"golang.org/x/time/rate"
)

type Consumer struct {
	fileService service.FileService
	reader      kafka.Consumer
	rateLimiter *rate.Limiter
	cfg         *config.ConsumerConfig
	wg          sync.WaitGroup
	sem         *semaphore.Weighted
}

func NewConsumer(cfg *config.ConsumerConfig, reader kafka.Consumer, fileService service.FileService) *Consumer {
	return &Consumer{
		reader:      reader,
		fileService: fileService,
		rateLimiter: rate.NewLimiter(
			rate.Limit(cfg.RateLimitPerSecond),
			cfg.RateLimitPerSecond+cfg.RateLimitPerSecond/100*20,
		),
		cfg: cfg,
		sem: semaphore.NewWeighted(int64(cfg.MaxInProcess)),
	}
}

func (r *Consumer) Close(_ context.Context) error {
	r.wg.Wait()
	return nil
}
