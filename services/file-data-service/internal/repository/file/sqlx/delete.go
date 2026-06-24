package sqlx

import (
	"context"
	"data/internal/repository"
	"database/sql"
	"errors"
)

func (r *fileRepository) DeleteFile(ctx context.Context, id string) error {
	res, err := r.db.DB.ExecContext(ctx, "DELETE FROM files WHERE id = $1", id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return repository.ErrNotFound
		}
		return err
	}
	rowsCount, err := res.RowsAffected()
	if err != nil {
		return err
	}
	if rowsCount == 0 {
		return repository.ErrNotFound
	}

	return nil
}
