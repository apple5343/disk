package redis

import (
	"context"
	"storage/internal/service"
	"storage/pkg/logger"

	"github.com/redis/go-redis/v9"
	"go.uber.org/zap"
)

const (
	uploadingCancelChannel = "uploading_cancel"
)

type Subscriber struct {
	rdb     *redis.Client
	service service.StorageService
}

func NewSubscriber(rdb *redis.Client, service service.StorageService) *Subscriber {
	return &Subscriber{
		rdb:     rdb,
		service: service,
	}
}

func (s *Subscriber) Listen(ctx context.Context) error {
	pubsub := s.rdb.Subscribe(ctx, uploadingCancelChannel)
	defer pubsub.Close()

	_, err := pubsub.Receive(ctx)
	if err != nil {
		return err
	}
	l, lOk := logger.FromContext(ctx)
	if !lOk {
		l = logger.NewBaseLogger()
	}

	ch := pubsub.Channel()
	for msg := range ch {
		err = s.service.HandleCancelUpload(ctx, msg.Payload)
		if err != nil {
			l.Error(ctx, "handle cancel uploading:"+err.Error(), zap.String("id", msg.Payload))
		}
	}
	return nil
}
