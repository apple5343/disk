package sqlx

import (
	"context"
	"data/internal/repository"
	"database/sql"
	"errors"
)

func (r *folderRepository) DeleteFolder(ctx context.Context, id string) error {
	_, err := r.db.DB.ExecContext(ctx, "DELETE FROM folders WHERE id = $1", id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return repository.ErrNotFound
		}
		return err
	}
	return nil
}
