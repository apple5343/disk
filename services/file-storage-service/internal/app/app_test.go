package app

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestApp(t *testing.T) {
	t.Parallel()
	app := NewApp()
	require.NotNil(t, app)
}
