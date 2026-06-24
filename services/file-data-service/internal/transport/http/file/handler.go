package file

import "data/internal/service"

type Handler struct {
	s service.FileService
}

func NewHandler(s service.FileService) *Handler {
	return &Handler{
		s: s,
	}
}
