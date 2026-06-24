package sqlx

import (
	"data/internal/repository"

	"github.com/jmoiron/sqlx"
)

type folderRepository struct {
	db *sqlx.DB
}

func NewRepository(db *sqlx.DB) repository.FolderRepository {
	return &folderRepository{db: db}
}
