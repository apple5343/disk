package models

import (
	"testing"

	"github.com/brianvoe/gofakeit"
	"github.com/stretchr/testify/require"
)

func TestFolder(t *testing.T) {
	t.Parallel()

	t.Run("valid folder", func(t *testing.T) {
		t.Parallel()
		parentID := gofakeit.UUID()
		folder := &Folder{
			UserID:   gofakeit.UUID(),
			Name:     gofakeit.Word(),
			ParentID: &parentID,
		}

		require.NoError(t, folder.BeforeCreate())
	})

	t.Run("/ in name", func(t *testing.T) {
		t.Parallel()
		folder := &Folder{
			Name: gofakeit.Word() + ".",
		}

		err := folder.BeforeCreate()
		require.Error(t, err)
	})

	t.Run("/ in name", func(t *testing.T) {
		t.Parallel()
		folder := &Folder{
			Name: "na/me",
		}

		err := folder.BeforeCreate()
		require.Error(t, err)
	})

	t.Run("* in name", func(t *testing.T) {
		t.Parallel()
		folder := &Folder{
			Name: gofakeit.Word() + "*",
		}

		err := folder.BeforeCreate()
		require.Error(t, err)
	})

	t.Run("empty name", func(t *testing.T) {
		t.Parallel()
		folder := &Folder{}

		err := folder.BeforeCreate()
		require.Error(t, err)
	})
}
