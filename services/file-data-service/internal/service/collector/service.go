package collector

import (
	"data/internal/service"

	"github.com/apple5343/errorx"
)

var (
	ErrInvalidToken   = errorx.NewError("invalid token", errorx.BadRequest)
	ErrDeleteRoot     = errorx.NewError("cannot delete root folder", errorx.BadRequest)
	ErrFolderNotEmpty = errorx.NewError("folder is not empty", errorx.BadRequest)
)

type collectorService struct {
	fileService   service.FileService
	folderService service.FolderService
}

func NewService(fileService service.FileService, folderService service.FolderService) service.CollectorService {
	return &collectorService{fileService: fileService, folderService: folderService}
}
