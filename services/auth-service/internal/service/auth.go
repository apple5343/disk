package service

import (
	"context"
	"errors"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"gitlab.crja72.ru/golang/2025/autumn/projects/go36/disk/services/auth-service/internal/client"
	"gitlab.crja72.ru/golang/2025/autumn/projects/go36/disk/services/auth-service/internal/model"
	"golang.org/x/crypto/bcrypt"
)

type AuthService struct {
	userClient client.UserClient
	jwtSecret  string
}

type AuthServiceInterface interface {
	Register(ctx context.Context, name, login, password string) error
	Login(ctx context.Context, login, password string) (string, error)
}

func NewAuthService(userClient client.UserClient, jwtSecret string) *AuthService {
	return &AuthService{
		userClient: userClient,
		jwtSecret:  jwtSecret,
	}
}

func (s *AuthService) Register(ctx context.Context, name, login, password string) error {
	userID := uuid.New()

	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}

	req := &model.UserCreateRequest{
		ID:           userID.String(),
		Login:        login,
		Name:         name,
		PasswordHash: hash,
	}

	return s.userClient.CreateUser(ctx, req)
}

func (s *AuthService) Login(ctx context.Context, login, password string) (string, error) {
	user, err := s.userClient.GetUserByLogin(ctx, login)
	if err != nil {
		if errors.Is(err, model.ErrUserNotFound) {
			return "", model.ErrInvalidCredentials
		}
		return "", err
	}

	if err = bcrypt.CompareHashAndPassword(user.PasswordHash, []byte(password)); err != nil {
		return "", model.ErrInvalidCredentials
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"UserID": user.ID.String(),
		"login":  user.Login,
		"exp":    time.Now().Add(1 * time.Hour).Unix(),
	})

	return token.SignedString([]byte(s.jwtSecret))
}
