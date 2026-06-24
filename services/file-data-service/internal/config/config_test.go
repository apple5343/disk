package config

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestConfig(t *testing.T) {
	_, err := NewPostgresConfig()
	require.NoError(t, err)
	_, err = NewLoggerConfig()
	require.NoError(t, err)
	_, err = NewJWTConfig()
	require.NoError(t, err)
	_, err = NewHTTPConfig()
	require.NoError(t, err)
	_, err = NewKafkaConsumerConfig()
	require.NoError(t, err)
	_, err = NewRedisConfig()
	require.NoError(t, err)
}
