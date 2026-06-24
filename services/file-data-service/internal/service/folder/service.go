package folder

import (
	"data/internal/repository"
	"data/internal/service"

	"github.com/apple5343/errorx"
)

var (
	ErrFolderNotFound = errorx.NewError("foder not found", errorx.BadRequest)
	ErrInvalidParent  = errorx.NewError("invalid parent", errorx.BadRequest)
	ErrFolderExists   = errorx.NewError("folder already exists", errorx.BadRequest)
	ErrInvalidID      = errorx.NewError("invalid id", errorx.BadRequest)
	ErrInvalidToken   = errorx.NewError("invalid token", errorx.BadRequest)
)

type folderService struct {
	folderRepository repository.FolderRepository
}

func NewService(folderRepository repository.FolderRepository) service.FolderService {
	return &folderService{folderRepository: folderRepository}
}
