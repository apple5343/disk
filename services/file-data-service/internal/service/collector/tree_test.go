package collector

import (
	"context"
	"data/internal/models"
	"data/internal/service/folder"
	"data/internal/service/mocks"
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

func randomFile() *models.FileMetadata {
	folder := "/" + gofakeit.Word() + "/" + gofakeit.Word()
	fileName := gofakeit.Word() + ".jpg"
	fullPath := folder + "/" + fileName
	return &models.FileMetadata{
		ID:          gofakeit.UUID(),
		UserID:      gofakeit.UUID(),
		StoragePath: "users" + fullPath,
		FolderID:    gofakeit.UUID(),
		FileName:    fileName,
		FullPath:    fullPath,
		Bucket:      "test-bucket",
		Size:        int64(gofakeit.Uint16()),
		ContentType: "image/jpeg",
		Tags: map[string]string{
			"tag1": "value1",
			"tag2": "value2",
		},
	}
}

func TestTree(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		parent := randomFolder()
		folders := []*models.Folder{randomFolder(), randomFolder()}
		files := []*models.FileMetadata{randomFile(), randomFile()}

		fileService := mocks.NewFileService(t)
		folderService := mocks.NewFolderService(t)

		fileService.EXPECT().
			GetFilesByFolderID(mock.Anything, mock.Anything).
			Return(files, nil).
			Run(func(_ context.Context, id string) {
				require.Equal(t, parent.ID, id)
			})

		folderService.EXPECT().
			GetFolderByID(mock.Anything, mock.Anything).
			Return(parent, nil).
			Run(func(_ context.Context, id string) {
				require.Equal(t, parent.ID, id)
			})

		folderService.EXPECT().
			GetFoldersByParentID(mock.Anything, mock.Anything).
			Return(folders, nil).
			Run(func(_ context.Context, id string) {
				require.Equal(t, parent.ID, id)
			})

		service := NewService(fileService, folderService)
		ctx := ctxutil.ContextWithUserID(context.Background(), parent.UserID)

		res, err := service.GetFolderTree(ctx, parent.ID)
		require.NoError(t, err)
		require.Len(t, res.Childs, len(folders))
		require.Len(t, res.Files, len(files))

		for i, folder := range folders {
			require.Equal(t, folder, res.Childs[i])
		}

		for i, file := range files {
			require.Equal(t, file, res.Files[i])
		}
	})

	t.Run("folder not found", func(t *testing.T) {
		folderID := gofakeit.UUID()
		folderService := mocks.NewFolderService(t)

		folderService.EXPECT().
			GetFolderByID(mock.Anything, mock.Anything).
			Return(nil, folder.ErrFolderNotFound).
			Run(func(_ context.Context, id string) {
				require.Equal(t, folderID, id)
			})

		service := NewService(nil, folderService)
		ctx := ctxutil.ContextWithUserID(context.Background(), gofakeit.UUID())

		_, err := service.GetFolderTree(ctx, folderID)
		require.ErrorIs(t, err, folder.ErrFolderNotFound)
	})
}
