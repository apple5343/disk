package service

import (
	"context"
	"errors"
	"testing"

	"github.com/google/uuid"
	"gitlab.crja72.ru/golang/2025/autumn/projects/go36/disk/services/user-service.v2/internal/model"
)

// fakeRepository — простая заглушка без mock-библиотек.
type fakeRepository struct {
	getByIDErr    error
	getByLoginErr error
	createErr     error
}

func (f *fakeRepository) GetByID(_ context.Context, id uuid.UUID) (*model.User, error) {
	if f.getByIDErr != nil {
		return nil, f.getByIDErr
	}
	return &model.User{ID: id, Login: "test", Name: "Test"}, nil
}

func (f *fakeRepository) GetByLogin(_ context.Context, login string) (*model.InternalUser, error) {
	if f.getByLoginErr != nil {
		return nil, f.getByLoginErr
	}
	return &model.InternalUser{Login: login}, nil
}

func (f *fakeRepository) Create(_ context.Context, _ *model.User) error {
	return f.createErr
}

func TestUserService_Create(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		repo := &fakeRepository{}
		svc := NewUserService(repo)

		user := &model.User{
			ID:           uuid.MustParse("123e4567-e89b-12d3-a456-426614174000"),
			Login:        "testuser",
			Name:         "Test User",
			PasswordHash: []byte("hash"),
		}

		err := svc.Create(context.Background(), user)
		if err != nil {
			t.Errorf("expected no error, got %v", err)
		}
	})

	t.Run("validation error", func(t *testing.T) {
		repo := &fakeRepository{}
		svc := NewUserService(repo)

		user := &model.User{
			ID:           uuid.MustParse("123e4567-e89b-12d3-a456-426614174000"),
			Login:        "",
			Name:         "Test",
			PasswordHash: []byte("hash"),
		}

		err := svc.Create(context.Background(), user)
		if err == nil {
			t.Error("expected validation error")
		}
	})

	t.Run("user exists", func(t *testing.T) {
		repo := &fakeRepository{
			createErr: errors.New("duplicate key value violates unique constraint"),
		}
		svc := NewUserService(repo)

		user := &model.User{
			ID:           uuid.MustParse("123e4567-e89b-12d3-a456-426614174000"),
			Login:        "testuser",
			Name:         "Test User",
			PasswordHash: []byte("hash"),
		}

		err := svc.Create(context.Background(), user)
		if !errors.Is(err, ErrUserExists) {
			t.Errorf("expected ErrUserExists, got %v", err)
		}
	})

	t.Run("unknown db error", func(t *testing.T) {
		repo := &fakeRepository{
			createErr: errors.New("connection timeout"),
		}
		svc := NewUserService(repo)

		user := &model.User{
			ID:           uuid.MustParse("123e4567-e89b-12d3-a456-426614174000"),
			Login:        "testuser",
			Name:         "Test User",
			PasswordHash: []byte("hash"),
		}

		err := svc.Create(context.Background(), user)
		if err == nil {
			t.Error("expected db error")
		}
		if errors.Is(err, ErrUserExists) {
			t.Error("should not be ErrUserExists")
		}
	})
}

func TestUserService_GetByID(t *testing.T) {
	repo := &fakeRepository{}
	svc := NewUserService(repo)

	_, err := svc.GetByID(context.Background(), uuid.Nil)
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}
}
