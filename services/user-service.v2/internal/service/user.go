package service

import (
	"context"
	"errors"
	"strings"

	"github.com/google/uuid"
	"gitlab.crja72.ru/golang/2025/autumn/projects/go36/disk/services/user-service.v2/internal/model"
	"gitlab.crja72.ru/golang/2025/autumn/projects/go36/disk/services/user-service.v2/internal/repository"
)

var ErrUserExists = errors.New("user already exists")

type UserService struct {
	repo repository.UserRepository
}

func NewUserService(repo repository.UserRepository) *UserService {
	return &UserService{repo: repo}
}

func (s *UserService) GetByID(ctx context.Context, id uuid.UUID) (*model.User, error) {
	return s.repo.GetByID(ctx, id)
}

func (s *UserService) GetByLogin(ctx context.Context, login string) (*model.InternalUser, error) {
	return s.repo.GetByLogin(ctx, login)
}

func (s *UserService) Create(ctx context.Context, user *model.User) error {
	if err := user.Validate(); err != nil {
		return err
	}
	err := s.repo.Create(ctx, user)
	if err != nil {
		if strings.Contains(err.Error(), "unique_constraint") ||
			strings.Contains(err.Error(), "duplicate key") ||
			strings.Contains(err.Error(), "UNIQUE constraint") {
			return ErrUserExists
		}
		return err
	}
	return nil
}
