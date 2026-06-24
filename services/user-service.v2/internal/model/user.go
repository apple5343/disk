package model

import (
	"errors"
	"strings"

	"github.com/google/uuid"
)

type User struct {
	ID           uuid.UUID `json:"id"    db:"id"`
	Login        string    `json:"login" db:"login"`
	Name         string    `json:"name"  db:"name"`
	PasswordHash []byte    `json:"-"     db:"password_hash"`
}

type InternalUser struct {
	ID           uuid.UUID `json:"id"`
	Login        string    `json:"login"`
	Name         string    `json:"name"`
	PasswordHash []byte    `json:"password_hash"`
}

var (
	ErrIDIsRequired = errors.New("id is required")
	ErrLoginLength  = errors.New("login must be 3-64 characters")
	ErrWrongName    = errors.New("name is required and max 128 chars")
	ErrPassReqired  = errors.New("password_hash is required")
)

func (u *User) Validate() error {
	if u.ID == uuid.Nil {
		return ErrIDIsRequired
	}
	if trimmed := strings.TrimSpace(u.Login); trimmed == "" || len(trimmed) < 3 || len(trimmed) > 64 {
		return ErrLoginLength
	}
	if trimmed := strings.TrimSpace(u.Name); trimmed == "" || len(trimmed) > 128 {
		return ErrWrongName
	}
	if len(u.PasswordHash) == 0 {
		return ErrPassReqired
	}
	return nil
}
