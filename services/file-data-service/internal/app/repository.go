package app

import (
	"data/internal/repository/file/redis"
	file "data/internal/repository/file/sqlx"
	folder "data/internal/repository/folder/sqlx"

	"go.uber.org/fx"
)

func RepositoryModule() fx.Option {
	return fx.Module(
		"repository",
		file.NewModule(),
		redis.NewModule(),
		folder.NewModule(),
	)
}
