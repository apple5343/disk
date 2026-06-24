package kafka

import (
	"context"
	"data/internal/infrastructure/kafka"
	"data/pkg/logger"
	"encoding/json"
	"errors"
)

const (
	UploadingStatus = "uploading"
	UploadedStatus  = "uploaded"
	FailedStatus    = "failed"
)

func (c *Consumer) Consume(ctx context.Context) error {
	for {
		select {
		case <-ctx.Done():
			c.wg.Wait()
			return nil
		default:
			if err := c.rateLimiter.Wait(ctx); err != nil {
				if errors.Is(err, context.Canceled) {
					return nil
				}
				return err
			}

			if err := c.sem.Acquire(ctx, 1); err != nil {
				if errors.Is(err, context.Canceled) {
					return nil
				}
				return err
			}

			msg, err := c.reader.Read(ctx)
			if err != nil {
				c.sem.Release(1)
				return err
			}

			c.wg.Add(1)
			go c.ProcessMessage(ctx, msg)
		}
	}
}

func (c *Consumer) ProcessMessage(ctx context.Context, msg kafka.Message) {
	defer c.wg.Done()
	defer c.sem.Release(1)

	l, lOk := logger.FromContext(ctx)
	if !lOk {
		l = logger.NewBaseLogger()
	}
	if err := c.HandleMessage(ctx, msg); err != nil {
		l.Error(ctx, "handle message: "+err.Error())
		return
	}

	if err := msg.Commit(ctx); err != nil {
		l.Error(ctx, "handle message: "+err.Error())
	}
}

func (c *Consumer) HandleMessage(ctx context.Context, msg kafka.Message) error {
	var message Message
	if err := json.Unmarshal(msg.Value(), &message); err != nil {
		return err
	}
	switch message.Status {
	case UploadingStatus:
		return c.fileService.FileUploading(ctx, FileMetadataFromKafka(message.File))
	case UploadedStatus:
		return c.fileService.FileUploaded(ctx, FileMetadataFromKafka(message.File))
	case FailedStatus:
		return c.fileService.FileFailed(ctx, FileMetadataFromKafka(message.File), nil)
	}
	return nil
}
