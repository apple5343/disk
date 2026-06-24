package handler

import (
	"errors"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"gitlab.crja72.ru/golang/2025/autumn/projects/go36/disk/services/user-service.v2/internal/model"
	"gitlab.crja72.ru/golang/2025/autumn/projects/go36/disk/services/user-service.v2/internal/repository"
	"gitlab.crja72.ru/golang/2025/autumn/projects/go36/disk/services/user-service.v2/internal/service"
)

type UserHandler struct {
	service *service.UserService
}

func NewUserHandler(service *service.UserService) *UserHandler {
	return &UserHandler{service: service}
}

type CreateUserRequest struct {
	ID           string `json:"id"`
	Login        string `json:"login"`
	Name         string `json:"name"`
	PasswordHash []byte `json:"password_hash"`
}

func (h *UserHandler) CreateUser(c *fiber.Ctx) error {
	var req CreateUserRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "invalid json",
		})
	}

	userID, err := uuid.Parse(req.ID)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "invalid user ID format",
		})
	}

	if len(req.PasswordHash) == 0 {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "password_hash is required",
		})
	}

	user := &model.User{
		ID:           userID,
		Login:        req.Login,
		Name:         req.Name,
		PasswordHash: req.PasswordHash,
	}

	if err = h.service.Create(c.Context(), user); err != nil {
		if errors.Is(err, service.ErrUserExists) {
			return c.Status(fiber.StatusConflict).JSON(fiber.Map{
				"error": "user already exists",
			})
		}
		if errors.Is(err, model.ErrLoginLength) ||
			errors.Is(err, model.ErrWrongName) ||
			errors.Is(err, model.ErrPassReqired) {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error": err.Error(),
			})
		}
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "internal error",
		})
	}

	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"id":    user.ID.String(),
		"login": user.Login,
		"name":  user.Name,
	})
}

func (h *UserHandler) GetUser(c *fiber.Ctx) error {
	idStr := c.Params("id")
	userID, err := uuid.Parse(idStr)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "invalid user ID",
		})
	}

	user, err := h.service.GetByID(c.Context(), userID)
	if err != nil {
		if errors.Is(err, repository.ErrUserNotFound) {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
				"error": "user not found",
			})
		}
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "internal error",
		})
	}

	return c.JSON(user)
}

func (h *UserHandler) GetUserByLogin(c *fiber.Ctx) error {
	login := c.Params("login")
	if login == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "login is required",
		})
	}

	user, err := h.service.GetByLogin(c.Context(), login)
	if err != nil {
		if errors.Is(err, repository.ErrUserNotFound) {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
				"error": "user not found",
			})
		}
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "internal error",
		})
	}

	return c.JSON(user)
}
