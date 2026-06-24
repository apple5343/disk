package kafka

import "context"

type Producer interface {
	Publish(ctx context.Context, topic string, key []byte, value []byte) error
	Close(ctx context.Context) error
	Ping(ctx context.Context) error
}
