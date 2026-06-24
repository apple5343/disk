package kafka

import "context"

type Consumer interface {
	Read(ctx context.Context) (Message, error)
	Ping(ctx context.Context) error
	Close(ctx context.Context) error
}

type Message interface {
	Value() []byte
	Commit(ctx context.Context) error
}
