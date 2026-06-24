package storage

import "storage/internal/service"

type Handler struct {
	s service.StorageService
}

func NewHandler(s service.StorageService) *Handler {
	return &Handler{
		s: s,
	}
}
