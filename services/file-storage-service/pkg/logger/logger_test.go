package logger

import (
	"context"
	"storage/internal/config"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestLogger(t *testing.T) {
	t.Parallel()
	l := NewLogger(&config.LoggerConfig{Level: "dev"})
	require.NotNil(t, l)
	l.Debug(context.Background(), "test")
	l.Info(context.Background(), "test")
	l.Error(context.Background(), "test")
}
