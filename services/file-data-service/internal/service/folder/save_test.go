package folder

import (
	"context"
	"data/internal/models"
	"data/internal/repository"
	"data/internal/repository/mocks"
	ctxutil "data/internal/utils/ctx"
	"testing"
	"time"

	"github.com/brianvoe/gofakeit"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

func TestSave(t *testing.T) {
	t.Parallel()

	t.Run("success without parent", func(t *testing.T) {
		t.Parallel()

		folder := randomFolder()
		root := &models.Folder{
			ID:        gofakeit.UUID(),
			UserID:    folder.UserID,
			Name:      "root",
			ParentID:  nil,
			FullPath:  "/",
			PathDepth: 0,
			IsRoot:    true,
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		}
		folder.ParentID = nil
		fr := mocks.NewFolderRepository(t)

		fr.EXPECT().
			GetRootFolder(mock.Anything, mock.Anything).
			Return(root, nil).
			Run(func(_ context.Context, userID string) {
				require.Equal(t, folder.UserID, userID)
			})
		fr.EXPECT().
			GetFolderByID(mock.Anything, mock.Anything).
			Return(root, nil).
			Run(func(_ context.Context, id string) {
				require.Equal(t, root.ID, id)
			})

		fr.EXPECT().
			SaveFolder(mock.Anything, mock.Anything).
			Return(folder, nil).
			Run(func(_ context.Context, f *models.Folder) {
				require.Equal(t, root.ID, *f.ParentID)
			})

		service := NewService(fr)
		ctx := ctxutil.ContextWithUserID(context.Background(), folder.UserID)

		res, err := service.SaveFolder(ctx, folder)
		require.NoError(t, err)
		require.Equal(t, "/"+folder.Name, res.FullPath)
		require.Equal(t, 1, res.PathDepth)
		require.Equal(t, root.ID, *res.ParentID)
	})

	t.Run("success with parent", func(t *testing.T) {
		t.Parallel()

		folder := randomFolder()
		parent := randomFolder()
		folder.ParentID = &parent.ID
		folder.UserID = parent.UserID
		fr := mocks.NewFolderRepository(t)
		fr.EXPECT().
			GetFolderByID(mock.Anything, mock.Anything).
			Return(parent, nil).
			Run(func(_ context.Context, id string) {
				require.Equal(t, parent.ID, id)
			})

		fr.EXPECT().
			SaveFolder(mock.Anything, mock.Anything).
			Return(folder, nil).
			Run(func(_ context.Context, f *models.Folder) {
				require.Equal(t, parent.ID, *f.ParentID)
			})

		service := NewService(fr)
		ctx := ctxutil.ContextWithUserID(context.Background(), folder.UserID)

		res, err := service.SaveFolder(ctx, folder)
		require.NoError(t, err)
		require.Equal(t, parent.FullPath+"/"+folder.Name, res.FullPath)
		require.Equal(t, parent.PathDepth+1, res.PathDepth)
		require.Equal(t, parent.ID, *res.ParentID)
	})

	t.Run("alredy exists", func(t *testing.T) {
		t.Parallel()

		folder := randomFolder()
		parent := randomFolder()
		folder.ParentID = &parent.ID
		folder.UserID = parent.UserID
		fr := mocks.NewFolderRepository(t)
		fr.EXPECT().
			GetFolderByID(mock.Anything, mock.Anything).
			Return(parent, nil).
			Run(func(_ context.Context, id string) {
				require.Equal(t, parent.ID, id)
			})

		fr.EXPECT().
			SaveFolder(mock.Anything, mock.Anything).
			Return(nil, repository.ErrAlredyExists).
			Run(func(_ context.Context, f *models.Folder) {
				require.Equal(t, parent.ID, *f.ParentID)
			})

		service := NewService(fr)
		ctx := ctxutil.ContextWithUserID(context.Background(), folder.UserID)

		_, err := service.SaveFolder(ctx, folder)
		require.ErrorIs(t, err, ErrFolderExists)
	})

	t.Run("invalid parent", func(t *testing.T) {
		t.Parallel()

		folder := randomFolder()
		parent := randomFolder()
		folder.ParentID = &parent.ID
		folder.UserID = parent.UserID
		fr := mocks.NewFolderRepository(t)
		fr.EXPECT().
			GetFolderByID(mock.Anything, mock.Anything).
			Return(nil, repository.ErrNotFound).
			Run(func(_ context.Context, id string) {
				require.Equal(t, parent.ID, id)
			})

		service := NewService(fr)
		ctx := ctxutil.ContextWithUserID(context.Background(), folder.UserID)

		_, err := service.SaveFolder(ctx, folder)
		require.ErrorIs(t, err, ErrInvalidParent)
	})

	t.Run("empty user id", func(t *testing.T) {
		service := NewService(nil)
		_, err := service.SaveFolder(context.Background(), nil)
		require.ErrorIs(t, err, ErrInvalidToken)
	})
}

func TestRootFolder(t *testing.T) {
	t.Parallel()

	t.Run("create", func(t *testing.T) {
		t.Parallel()

		userID := gofakeit.UUID()
		root := &models.Folder{
			UserID:    userID,
			Name:      "root",
			ParentID:  nil,
			FullPath:  "/",
			PathDepth: 0,
			IsRoot:    true,
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		}
		fr := mocks.NewFolderRepository(t)
		fr.EXPECT().
			GetRootFolder(mock.Anything, mock.Anything).
			Return(nil, repository.ErrNotFound).
			Run(func(_ context.Context, id string) {
				require.Equal(t, userID, id)
			})
		fr.EXPECT().
			SaveFolder(mock.Anything, mock.Anything).
			Return(root, nil).
			Run(func(_ context.Context, f *models.Folder) {
				require.Equal(t, userID, f.UserID)
				require.Equal(t, "root", f.Name)
				require.Nil(t, f.ParentID)
				require.Equal(t, "/", f.FullPath)
				require.Equal(t, 0, f.PathDepth)
				require.True(t, f.IsRoot)
			})

		service := NewService(fr)
		ctx := ctxutil.ContextWithUserID(context.Background(), userID)

		res, err := service.RootFolder(ctx, userID)
		require.NoError(t, err)
		require.Equal(t, root, res)
	})
}
