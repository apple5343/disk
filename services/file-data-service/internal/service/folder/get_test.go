package folder

import (
	"context"
	"data/internal/models"
	"data/internal/repository"
	"data/internal/repository/mocks"
	ctxutil "data/internal/utils/ctx"
	"testing"

	"github.com/brianvoe/gofakeit"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

func randomFolder() *models.Folder {
	folderName := gofakeit.Word()
	parentID := gofakeit.UUID()
	return &models.Folder{
		ID:        gofakeit.UUID(),
		UserID:    gofakeit.UUID(),
		ParentID:  &parentID,
		Name:      folderName,
		FullPath:  "/" + gofakeit.Word() + "/" + folderName,
		PathDepth: 2,
	}
}

func TestGet(t *testing.T) {
	t.Parallel()

	t.Run("success", func(t *testing.T) {
		t.Parallel()

		folder := randomFolder()
		fr := mocks.NewFolderRepository(t)

		fr.EXPECT().
			GetFolderByID(mock.Anything, mock.Anything).
			Return(folder, nil).
			Run(func(_ context.Context, id string) {
				require.Equal(t, folder.ID, id)
			})

		service := NewService(fr)
		ctx := ctxutil.ContextWithUserID(context.Background(), folder.UserID)

		res, err := service.GetFolderByID(ctx, folder.ID)
		require.NoError(t, err)
		require.Equal(t, *folder, *res)
	})

	t.Run("folder not found", func(t *testing.T) {
		t.Parallel()

		folder := randomFolder()
		fr := mocks.NewFolderRepository(t)

		fr.EXPECT().
			GetFolderByID(mock.Anything, mock.Anything).
			Return(nil, repository.ErrNotFound).
			Run(func(_ context.Context, id string) {
				require.Equal(t, folder.ID, id)
			})

		service := NewService(fr)
		ctx := ctxutil.ContextWithUserID(context.Background(), folder.UserID)

		_, err := service.GetFolderByID(ctx, folder.ID)
		require.ErrorIs(t, err, ErrFolderNotFound)
	})

	t.Run("user id does not match", func(t *testing.T) {
		t.Parallel()

		folder := randomFolder()
		userID := gofakeit.UUID()
		fr := mocks.NewFolderRepository(t)

		fr.EXPECT().
			GetFolderByID(mock.Anything, mock.Anything).
			Return(folder, nil).
			Run(func(_ context.Context, id string) {
				require.Equal(t, folder.ID, id)
			})

		service := NewService(fr)
		ctx := ctxutil.ContextWithUserID(context.Background(), userID)

		_, err := service.GetFolderByID(ctx, folder.ID)
		require.ErrorIs(t, err, ErrFolderNotFound)
	})

	t.Run("empty user id", func(t *testing.T) {
		t.Parallel()
		service := NewService(nil)
		_, err := service.GetFolderByID(context.Background(), gofakeit.UUID())
		require.ErrorIs(t, err, ErrInvalidToken)
	})

	t.Run("invalid id", func(t *testing.T) {
		t.Parallel()
		service := NewService(nil)
		_, err := service.GetFolderByID(ctxutil.ContextWithUserID(context.Background(), gofakeit.UUID()), "invalid-id")
		require.ErrorIs(t, err, ErrInvalidID)
	})
}

func TestGetByParentID(t *testing.T) {
	t.Parallel()

	t.Run("success", func(t *testing.T) {
		t.Parallel()

		folders := []*models.Folder{randomFolder(), randomFolder()}
		parentID := *folders[0].ParentID
		fr := mocks.NewFolderRepository(t)

		fr.EXPECT().
			GetFolderByParentID(mock.Anything, mock.Anything).
			Return(folders, nil).
			Run(func(_ context.Context, id string) {
				require.Equal(t, parentID, id)
			})

		service := NewService(fr)
		ctx := ctxutil.ContextWithUserID(context.Background(), folders[0].UserID)

		res, err := service.GetFoldersByParentID(ctx, parentID)
		require.NoError(t, err)
		require.Len(t, res, len(folders))
	})

	t.Run("invalid parent id", func(t *testing.T) {
		service := NewService(nil)
		ctx := ctxutil.ContextWithUserID(context.Background(), gofakeit.UUID())
		_, err := service.GetFoldersByParentID(ctx, "invalid-id")
		require.ErrorIs(t, err, ErrInvalidID)
	})
}
