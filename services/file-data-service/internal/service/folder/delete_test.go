package folder

import (
	"context"
	"data/internal/repository/mocks"
	ctxutil "data/internal/utils/ctx"
	"testing"

	"github.com/brianvoe/gofakeit"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

func TestDelete(t *testing.T) {
	t.Parallel()

	t.Run("success", func(t *testing.T) {
		t.Parallel()
		fr := mocks.NewFolderRepository(t)
		folderID := gofakeit.UUID()

		fr.EXPECT().
			DeleteFolder(mock.Anything, mock.Anything).
			Return(nil).
			Run(func(_ context.Context, id string) {
				require.Equal(t, folderID, id)
			})

		service := NewService(fr)
		ctx := ctxutil.ContextWithUserID(context.Background(), gofakeit.UUID())

		err := service.DeleteFolder(ctx, folderID)
		require.NoError(t, err)
	})
}
