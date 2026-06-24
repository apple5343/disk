package file

import (
	"data/internal/repository"
	"data/internal/service"

	"github.com/apple5343/errorx"
)

var (
	ErrFileNotFound   = errorx.NewError("file not found", errorx.BadRequest)
	ErrInvalidToken   = errorx.NewError("invalid token", errorx.Unauthorized)
	ErrFolderNotFound = errorx.NewError("folder not found", errorx.BadRequest)
	ErrInvalidID      = errorx.NewError("invalid id", errorx.BadRequest)
)

type fileService struct {
	fileRepository  repository.FileRepository
	folderService   service.FolderService
	fileProccessing repository.FileProcessingRepository
}

func NewService(
	repo repository.FileRepository,
	folderService service.FolderService,
	fileProccessing repository.FileProcessingRepository,
) service.FileService {
	return &fileService{fileRepository: repo, folderService: folderService, fileProccessing: fileProccessing}
}
