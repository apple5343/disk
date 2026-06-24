package app

import (
	"storage/internal/repository/storage"
	"storage/internal/repository/upload"

	"go.uber.org/fx"
)

func RepositoryModule() fx.Option {
	return fx.Module(
		"repository",
		upload.NewModule(),
		storage.NewModule(),
	)
}
