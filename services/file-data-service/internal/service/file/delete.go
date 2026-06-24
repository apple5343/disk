package file

import (
	"context"
	ctxutil "data/internal/utils/ctx"

	"github.com/google/uuid"
)

func (s *fileService) DeleteFile(ctx context.Context, id string) error {
	userID := ctxutil.UserIDFromContext(ctx)
	if userID == "" {
		return ErrInvalidToken
	}

	_, err := uuid.Parse(id)
	if err != nil {
		return ErrInvalidID
	}

	_, err = s.GetFileByID(ctx, id)
	if err != nil {
		return err
	}

	return s.fileRepository.DeleteFile(ctx, id)
}
