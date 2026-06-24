package app

import (
	"storage/internal/adapter/file"
	"storage/internal/adapter/upload"
	"storage/internal/service/storage"

	"go.uber.org/fx"
)

func ServiceModule() fx.Option {
	return fx.Module(
		"service",
		file.NewModule(),
		upload.NewModule(),
		storage.NewModule(),
	)
}
