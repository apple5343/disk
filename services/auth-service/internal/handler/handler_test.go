package handler

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gofiber/fiber/v2"
	"gitlab.crja72.ru/golang/2025/autumn/projects/go36/disk/services/auth-service/internal/model"
)

type fakeAuthService struct {
	registerErr error
	loginErr    error
}

func (f *fakeAuthService) Register(_ context.Context, _, _, _ string) error {
	return f.registerErr
}

func (f *fakeAuthService) Login(_ context.Context, _, _ string) (string, error) {
	if f.loginErr != nil {
		return "", f.loginErr
	}
	return "fake-token", nil
}

func TestAuthHandler_Register(t *testing.T) {
	app := fiber.New()

	// ok
	service1 := &fakeAuthService{}
	handler1 := NewAuthHandler(service1)
	app.Post("/register", handler1.Register)

	reqBody, _ := json.Marshal(RegisterRequest{
		Name:     "Test User",
		Login:    "testuser",
		Password: "password123",
	})
	req := httptest.NewRequest(http.MethodPost, "/register", bytes.NewReader(reqBody))
	req.Header.Set("Content-Type", "application/json")
	resp, _ := app.Test(req)
	if resp.StatusCode != http.StatusCreated {
		t.Errorf("expected 201, got %d", resp.StatusCode)
	}

	// fail 409
	service2 := &fakeAuthService{
		registerErr: model.ErrUserExists,
	}
	handler2 := NewAuthHandler(service2)
	app.Post("/register2", handler2.Register)

	req = httptest.NewRequest(http.MethodPost, "/register2", bytes.NewReader(reqBody))
	req.Header.Set("Content-Type", "application/json")
	resp, _ = app.Test(req)
	if resp.StatusCode != http.StatusConflict {
		t.Errorf("expected 409 for user exists, got %d", resp.StatusCode)
	}
}

func TestAuthHandler_Login(t *testing.T) {
	app := fiber.New()

	// ok
	service1 := &fakeAuthService{}
	handler1 := NewAuthHandler(service1)
	app.Post("/login", handler1.Login)

	reqBody, _ := json.Marshal(LoginRequest{
		Login:    "testuser",
		Password: "password123",
	})
	req := httptest.NewRequest(http.MethodPost, "/login", bytes.NewReader(reqBody))
	req.Header.Set("Content-Type", "application/json")
	resp, _ := app.Test(req)
	if resp.StatusCode != http.StatusOK {
		t.Errorf("expected 200, got %d", resp.StatusCode)
	}

	// fail 401
	service2 := &fakeAuthService{
		loginErr: model.ErrInvalidCredentials,
	}
	handler2 := NewAuthHandler(service2)
	app.Post("/login2", handler2.Login)

	req = httptest.NewRequest(http.MethodPost, "/login2", bytes.NewReader(reqBody))
	req.Header.Set("Content-Type", "application/json")
	resp, _ = app.Test(req)
	if resp.StatusCode != http.StatusUnauthorized {
		t.Errorf("expected 401 for invalid credentials, got %d", resp.StatusCode)
	}
}

func TestAuthHandler_InvalidJSON(t *testing.T) {
	app := fiber.New()
	service := &fakeAuthService{}
	handler := NewAuthHandler(service)
	app.Post("/register", handler.Register)

	req := httptest.NewRequest(http.MethodPost, "/register", bytes.NewReader([]byte("{invalid json}")))
	req.Header.Set("Content-Type", "application/json")
	resp, _ := app.Test(req)
	if resp.StatusCode != http.StatusBadRequest {
		t.Errorf("expected 400 for invalid JSON, got %d", resp.StatusCode)
	}
}
