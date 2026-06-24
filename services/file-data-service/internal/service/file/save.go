package file

import (
	"context"
	"data/internal/models"
	"data/internal/repository"
	"errors"

	"github.com/apple5343/errorx"
)

func (s *fileService) SaveFile(ctx context.Context, metadata *models.FileMetadata) (*models.FileMetadata, error) {
	if metadata.FolderID == "" {
		rootFolder, err := s.folderService.RootFolder(ctx, metadata.UserID)
		if err != nil {
			return nil, errorx.NewError("save file: "+err.Error(), errorx.Internal)
		}
		metadata.FolderID = rootFolder.ID
	}

	file, err := s.fileRepository.SaveFile(ctx, metadata)
	if nil == err {
		return file, nil
	}

	if errors.Is(err, repository.ErrAlredyExists) {
		file, err = s.fileRepository.GetFileByStoragePath(ctx, metadata.StoragePath)
		if err != nil {
			return nil, errorx.NewError("save file: "+err.Error(), errorx.Internal)
		}
		metadata.ID = file.ID
		file, err = s.fileRepository.UpdateFileByPath(ctx, metadata)
		if err != nil {
			return nil, errorx.NewError("save file: "+err.Error(), errorx.Internal)
		}
		return file, nil
	}
	return nil, errorx.NewError("save file: "+err.Error(), errorx.Internal)
}
