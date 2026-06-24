package models

import (
	"testing"

	"github.com/brianvoe/gofakeit"
	"github.com/stretchr/testify/require"
)

func TestFile(t *testing.T) {
	t.Parallel()

	t.Run("valid file", func(t *testing.T) {
		t.Parallel()
		file := &FileMetadata{
			UserID:      gofakeit.UUID(),
			FullPath:    "/files/file.txt",
			FolderID:    gofakeit.UUID(),
			FileName:    "file.txt",
			Bucket:      "users",
			ContentType: "text/plain",
		}

		require.NoError(t, file.BeforeCreate())
	})

	t.Run("invalid file", func(t *testing.T) {
		t.Parallel()
		file := &FileMetadata{}
		err := file.BeforeCreate()
		require.Error(t, err)
	})
}
