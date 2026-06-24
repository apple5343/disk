package config

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestConfig(t *testing.T) {
	t.Parallel()
	_, err := NewMinioConfig()
	require.NoError(t, err)
	_, err = NewLoggerConfig()
	require.NoError(t, err)
	_, err = NewJWTConfig()
	require.NoError(t, err)
	_, err = NewHTTPConfig()
	require.NoError(t, err)
	_, err = NewFileDataClientConfig()
	require.NoError(t, err)
}
