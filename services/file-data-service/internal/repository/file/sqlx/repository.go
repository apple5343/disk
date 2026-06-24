package sqlx

import (
	"data/internal/repository"

	"github.com/jmoiron/sqlx"
)

type fileRepository struct {
	db *sqlx.DB
}

func NewRepository(db *sqlx.DB) repository.FileRepository {
	return &fileRepository{db: db}
}
