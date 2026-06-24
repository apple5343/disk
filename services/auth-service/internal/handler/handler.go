package handler

import (
	"context"
	"errors"
	"net/http"

	"github.com/gofiber/fiber/v2"
	"gitlab.crja72.ru/golang/2025/autumn/projects/go36/disk/services/auth-service/internal/model"
	"gitlab.crja72.ru/golang/2025/autumn/projects/go36/disk/services/auth-service/internal/service"
)

const expiresIn = 3600

type AuthServiceInterface interface {
	Register(ctx context.Context, name, login, password string) error
	Login(ctx context.Context, login, password string) (string, error)
}

var _ AuthServiceInterface = &service.AuthService{}

type AuthHandler struct {
	service service.AuthServiceInterface // ← интерфейс, а не конкретная структура
}

func NewAuthHandler(service AuthServiceInterface) *AuthHandler {
	return &AuthHandler{service: service}
}

// type AuthHandler struct {
// 	service *service.AuthService
// }

// func NewAuthHandler(service *service.AuthService) *AuthHandler {
// 	return &AuthHandler{service: service}
// }

type RegisterRequest struct {
	Name     string `json:"name"`
	Login    string `json:"login"`
	Password string `json:"password"`
}

type LoginRequest struct {
	Login    string `json:"login"`
	Password string `json:"password"`
}

func (h *AuthHandler) Register(c *fiber.Ctx) error {
	var req RegisterRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{"error": "invalid json"})
	}

	if err := h.service.Register(c.Context(), req.Name, req.Login, req.Password); err != nil {
		if errors.Is(err, model.ErrUserExists) {
			return c.Status(http.StatusConflict).JSON(fiber.Map{"error": "user already exists"})
		}
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{"error": "registration failed"})
	}

	return c.Status(http.StatusCreated).JSON(fiber.Map{"message": "user created"})
}

func (h *AuthHandler) Login(c *fiber.Ctx) error {
	var req LoginRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{"error": "invalid json"})
	}

	token, err := h.service.Login(c.Context(), req.Login, req.Password)
	if err != nil {
		if errors.Is(err, model.ErrInvalidCredentials) {
			return c.Status(http.StatusUnauthorized).JSON(fiber.Map{"error": "invalid credentials"})
		}
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{"error": "login failed"})
	}

	return c.JSON(fiber.Map{
		"token":      token,
		"expires_in": expiresIn,
	})
}
