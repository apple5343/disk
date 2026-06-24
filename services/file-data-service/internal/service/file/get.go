package file

import (
	"context"
	"data/internal/models"
	"data/internal/repository"
	ctxutil "data/internal/utils/ctx"
	"errors"

	"github.com/apple5343/errorx"
	"github.com/google/uuid"
)

func (s *fileService) GetFileByID(ctx context.Context, id string) (*models.FileMetadata, error) {
	userID := ctxutil.UserIDFromContext(ctx)
	if userID == "" {
		return nil, ErrInvalidToken
	}
	_, err := uuid.Parse(id)
	if err != nil {
		return nil, ErrInvalidID
	}
	file, err := s.fileRepository.GetFileByID(ctx, id)
	if err != nil {
		if errors.Is(err, repository.ErrNotFound) {
			return nil, ErrFileNotFound
		}
		return nil, errorx.NewError("get file: "+err.Error(), errorx.Internal)
	}
	if file.UserID != userID {
		return nil, ErrFileNotFound
	}
	return file, nil
}

func (s *fileService) GetFileByStoragePath(ctx context.Context, storagePath string) (*models.FileMetadata, error) {
	file, err := s.fileRepository.GetFileByStoragePath(ctx, storagePath)
	if err != nil {
		if errors.Is(err, repository.ErrNotFound) {
			return nil, ErrFileNotFound
		}
		return nil, errorx.NewError("get file: "+err.Error(), errorx.Internal)
	}
	return file, nil
}

func (s *fileService) GetFilesByFolderID(ctx context.Context, folderID string) ([]*models.FileMetadata, error) {
	userID := ctxutil.UserIDFromContext(ctx)
	if userID == "" {
		return nil, ErrInvalidToken
	}
	_, err := uuid.Parse(folderID)
	if err != nil {
		return nil, ErrInvalidID
	}
	_, err = s.folderService.GetFolderByID(ctx, folderID)
	if err != nil {
		return nil, err
	}
	files, err := s.fileRepository.GetFilesByFolderID(ctx, folderID)
	if err != nil {
		return nil, errorx.NewError("get file: "+err.Error(), errorx.Internal)
	}
	return files, nil
}

func (s *fileService) SearchFiles(ctx context.Context, req *models.SearchRequest) ([]*models.FileMetadata, error) {
	userID := ctxutil.UserIDFromContext(ctx)
	if userID == "" {
		return nil, ErrInvalidToken
	}
	_, err := uuid.Parse(req.FolderID)
	if err != nil && req.FolderID != "" {
		return nil, ErrInvalidID
	}
	files, err := s.fileRepository.SearchFiles(ctx, userID, req)
	if err != nil {
		return nil, errorx.NewError("search file: "+err.Error(), errorx.Internal)
	}
	return files, nil
}
