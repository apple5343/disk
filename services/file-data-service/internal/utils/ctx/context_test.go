package ctxutil

import (
	"context"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestContext(t *testing.T) {
	t.Parallel()
	ctx := ContextWithUserID(context.Background(), "user-id")
	require.NotNil(t, ctx)
	userID := ctx.Value(UserIDCtxKey)
	id, ok := userID.(string)
	require.True(t, ok)
	require.Equal(t, "user-id", id)
}
