package folder

import "data/internal/service"

type Handler struct {
	s         service.FolderService
	collector service.CollectorService
}

func NewHandler(s service.FolderService, collector service.CollectorService) *Handler {
	return &Handler{
		s:         s,
		collector: collector,
	}
}
