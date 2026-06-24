package logger

import (
	"context"
	"storage/internal/config"

	"go.uber.org/zap"
)

type ctxKey string

const (
	loggerLoggerKey    ctxKey = "logger"
	loggerRequestIDKey ctxKey = "x-request-id"
	loggerTraceIDKey   ctxKey = "x-trace-id"
)

type Logger interface {
	Debug(ctx context.Context, msg string, fields ...zap.Field)
	Info(ctx context.Context, msg string, fields ...zap.Field)
	Error(ctx context.Context, msg string, fields ...zap.Field)
}

type logger struct {
	z *zap.Logger
}

func NewLogger(cfg *config.LoggerConfig) Logger {
	loggerCfg := zap.NewProductionConfig()
	loggerCfg.DisableStacktrace = true
	if cfg.Level == "dev" {
		loggerCfg.Level = zap.NewAtomicLevelAt(zap.DebugLevel)
	}

	z, _ := loggerCfg.Build()
	return &logger{
		z: z,
	}
}

func (l *logger) Debug(ctx context.Context, msg string, fields ...zap.Field) {
	if val, ok := RequestIDFromContext(ctx); ok {
		fields = append(fields, zap.String(string(loggerRequestIDKey), val))
	}

	l.z.Debug(msg, fields...)
}

func (l *logger) Info(ctx context.Context, msg string, fields ...zap.Field) {
	if val, ok := RequestIDFromContext(ctx); ok {
		fields = append(fields, zap.String(string(loggerRequestIDKey), val))
	}

	l.z.Info(msg, fields...)
}

func (l *logger) Error(ctx context.Context, msg string, fields ...zap.Field) {
	if val, ok := RequestIDFromContext(ctx); ok {
		fields = append(fields, zap.String(string(loggerRequestIDKey), val))
	}

	l.z.Error(msg, fields...)
}

func ContextWithRequestID(ctx context.Context, requestID string) context.Context {
	return context.WithValue(ctx, loggerRequestIDKey, requestID)
}

func ContextWithTraceID(ctx context.Context, traceID string) context.Context {
	return context.WithValue(ctx, loggerTraceIDKey, traceID)
}

func RequestIDFromContext(ctx context.Context) (string, bool) {
	val := ctx.Value(loggerRequestIDKey)
	if id, ok := val.(string); ok {
		return id, ok
	}

	return "", false
}

func ContextWithLogger(ctx context.Context, l Logger) context.Context {
	return context.WithValue(ctx, loggerLoggerKey, l)
}

func FromContext(ctx context.Context) (Logger, bool) {
	val := ctx.Value(loggerLoggerKey)
	if l, ok := val.(Logger); ok {
		return l, ok
	}

	return nil, false
}
