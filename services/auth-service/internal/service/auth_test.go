package service

import (
	"context"
	"errors"
	"testing"

	"github.com/google/uuid"
	"gitlab.crja72.ru/golang/2025/autumn/projects/go36/disk/services/auth-service/internal/model"
	"golang.org/x/crypto/bcrypt"
)

type mockUserClient struct {
	createErr     error
	getByLoginErr error
	user          *model.InternalUser
}

func (m *mockUserClient) CreateUser(_ context.Context, _ *model.UserCreateRequest) error {
	return m.createErr
}

func (m *mockUserClient) GetUserByLogin(_ context.Context, _ string) (*model.InternalUser, error) {
	if m.getByLoginErr != nil {
		return nil, m.getByLoginErr
	}
	return m.user, nil
}

func TestAuthService_Register(t *testing.T) {
	client := &mockUserClient{}
	svc := NewAuthService(client, "test-secret")

	err := svc.Register(context.Background(), "Test User", "testuser", "password123")
	if err != nil {
		t.Errorf("Register failed: %v", err)
	}
}

func TestAuthService_Register_UserExists(t *testing.T) {
	client := &mockUserClient{
		createErr: model.ErrUserExists,
	}
	svc := NewAuthService(client, "test-secret")

	err := svc.Register(context.Background(), "Test", "test", "pass")
	if !errors.Is(err, model.ErrUserExists) {
		t.Errorf("Expected ErrUserExists, got %v", err)
	}
}

func TestAuthService_Login_Success(t *testing.T) {
	hash, _ := bcrypt.GenerateFromPassword([]byte("password123"), bcrypt.DefaultCost)

	client := &mockUserClient{
		user: &model.InternalUser{
			ID:           uuid.MustParse("123e4567-e89b-12d3-a456-426614174000"),
			Login:        "testuser",
			PasswordHash: hash,
		},
	}
	svc := NewAuthService(client, "test-secret")

	token, err := svc.Login(context.Background(), "testuser", "password123")
	if err != nil {
		t.Fatalf("Login failed: %v", err)
	}
	if token == "" {
		t.Error("Expected non-empty token")
	}
}

func TestAuthService_Login_UserNotFound(t *testing.T) {
	client := &mockUserClient{
		getByLoginErr: model.ErrUserNotFound,
	}
	svc := NewAuthService(client, "test-secret")

	token, err := svc.Login(context.Background(), "nonexistent", "pass")
	if err == nil {
		t.Error("Expected error for non-existent user")
	}
	if token != "" {
		t.Error("Expected empty token")
	}
	if !errors.Is(err, model.ErrInvalidCredentials) {
		t.Errorf("Expected ErrInvalidCredentials, got %v", err)
	}
}

func TestAuthService_Login_WrongPassword(t *testing.T) {
	hash, _ := bcrypt.GenerateFromPassword([]byte("correctpass"), bcrypt.DefaultCost)

	client := &mockUserClient{
		user: &model.InternalUser{
			ID:           uuid.MustParse("123e4567-e89b-12d3-a456-426614174000"),
			PasswordHash: hash,
		},
	}
	svc := NewAuthService(client, "test-secret")

	token, err := svc.Login(context.Background(), "testuser", "wrongpass")
	if err == nil {
		t.Error("Expected error for wrong password")
	}
	if token != "" {
		t.Error("Expected empty token")
	}
	if !errors.Is(err, model.ErrInvalidCredentials) {
		t.Errorf("Expected ErrInvalidCredentials, got %v", err)
	}
}
