package app

import (
	"data/internal/service/collector"
	"data/internal/service/file"
	"data/internal/service/folder"

	"go.uber.org/fx"
)

func ServiceModule() fx.Option {
	return fx.Module(
		"service",
		file.NewModule(),
		folder.NewModule(),
		collector.NewModule(),
	)
}
