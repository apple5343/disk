package model

import (
	"errors"

	"github.com/google/uuid"
)

type UserCreateRequest struct {
	ID           string `json:"id"`
	Login        string `json:"login"`
	Name         string `json:"name"`
	PasswordHash []byte `json:"password_hash"`
}

type InternalUser struct {
	ID           uuid.UUID `json:"id"`
	Login        string    `json:"login"`
	Name         string    `json:"name"`
	PasswordHash []byte    `json:"password_hash"`
}

var (
	ErrUserExists         = errors.New("user already exists")
	ErrUserNotFound       = errors.New("user not found")
	ErrInvalidCredentials = errors.New("invalid credentials")
)
