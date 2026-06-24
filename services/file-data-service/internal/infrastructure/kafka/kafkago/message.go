package kafkago

import (
	"context"

	infrastructure "data/internal/infrastructure/kafka"

	"github.com/segmentio/kafka-go"
)

type Message struct {
	msg *kafka.Message

	commit func(ctx context.Context) error
}

func (c *consumer) NewMessage(msg *kafka.Message) infrastructure.Message {
	return &Message{
		msg: msg,
		commit: func(ctx context.Context) error {
			return c.reader.CommitMessages(ctx, *msg)
		},
	}
}

func (m *Message) Value() []byte {
	return m.msg.Value
}

func (m *Message) Commit(ctx context.Context) error {
	return m.commit(ctx)
}
