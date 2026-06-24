package logger

import (
	"context"

	"go.uber.org/zap"
)

type baseLogger struct {
	l *zap.Logger
}

func NewBaseLogger() Logger {
	return baseLogger{
		l: zap.NewNop(),
	}
}

func (b baseLogger) Error(_ context.Context, msg string, fields ...zap.Field) {
	b.l.Error(msg, fields...)
}

func (b baseLogger) Info(_ context.Context, msg string, fields ...zap.Field) {
	b.l.Info(msg, fields...)
}

func (b baseLogger) Debug(_ context.Context, msg string, fields ...zap.Field) {
	b.l.Debug(msg, fields...)
}
