package folder

import (
	"context"
	"data/internal/models"
	"data/internal/repository"
	ctxutil "data/internal/utils/ctx"
	"errors"
	"time"

	"github.com/apple5343/errorx"
)

func (s *folderService) SaveFolder(ctx context.Context, folder *models.Folder) (*models.Folder, error) {
	userID := ctxutil.UserIDFromContext(ctx)
	if userID == "" {
		return nil, ErrInvalidToken
	}

	folder.UserID = userID
	var parentID string
	if folder.ParentID != nil {
		parentID = *folder.ParentID
	} else {
		rootFolder, err := s.RootFolder(ctx, userID)
		if err != nil {
			return nil, err
		}
		parentID = rootFolder.ID
	}

	folder.ParentID = &parentID
	err := folder.BeforeCreate()
	if err != nil {
		return nil, errorx.NewError(err.Error(), errorx.BadRequest)
	}

	parentFolder, err := s.GetFolderByID(ctx, *folder.ParentID)
	if err != nil {
		if errors.Is(err, ErrFolderNotFound) {
			return nil, ErrInvalidParent
		}
		return nil, err
	}

	folder.PathDepth = parentFolder.PathDepth + 1
	folder.FullPath = parentFolder.FullPath + "/" + folder.Name
	if folder.PathDepth == 1 {
		folder.FullPath = parentFolder.FullPath + folder.Name
	}
	folder.IsRoot = false
	folder, err = s.folderRepository.SaveFolder(ctx, folder)
	if err != nil {
		if errors.Is(err, repository.ErrAlredyExists) {
			return nil, ErrFolderExists
		}
		return nil, err
	}
	return folder, nil
}

func (s *folderService) RootFolder(ctx context.Context, userID string) (*models.Folder, error) {
	folder, err := s.folderRepository.GetRootFolder(ctx, userID)
	if err != nil {
		if errors.Is(err, repository.ErrNotFound) {
			folder, err = s.CreateRootFolder(ctx, userID)
			if err != nil {
				return nil, err
			}
		} else {
			return nil, err
		}
	}
	return folder, nil
}

func (s *folderService) CreateRootFolder(ctx context.Context, userID string) (*models.Folder, error) {
	folder := &models.Folder{
		UserID:    userID,
		Name:      "root",
		ParentID:  nil,
		FullPath:  "/",
		PathDepth: 0,
		IsRoot:    true,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}
	return s.folderRepository.SaveFolder(ctx, folder)
}
