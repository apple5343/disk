package model

import (
	"testing"

	"github.com/google/uuid"
)

func TestUserValidate_Success(t *testing.T) {
	user := &User{
		ID:           uuid.MustParse("123e4567-e89b-12d3-a456-426614174000"),
		Login:        "testuser",
		Name:         "Test User",
		PasswordHash: []byte("hash"),
	}

	err := user.Validate()
	if err != nil {
		t.Errorf("Validate() failed: %v", err)
	}
}

func TestUserValidate_InvalidID(t *testing.T) {
	user := &User{
		ID:           uuid.Nil,
		Login:        "testuser",
		Name:         "Test User",
		PasswordHash: []byte("hash"),
	}

	err := user.Validate()
	if err == nil {
		t.Error("Expected validation error for nil ID")
	}
}

func TestUserValidate_InvalidLogin(t *testing.T) {
	cases := []string{"", "ab",
		"a12345678901234567890123456789012345678901234567890123456789012345"}
	for _, login := range cases {
		user := &User{
			ID:           uuid.MustParse("123e4567-e89b-12d3-a456-426614174000"),
			Login:        login,
			Name:         "Test User",
			PasswordHash: []byte("hash"),
		}
		err := user.Validate()
		if err == nil {
			t.Errorf("Expected validation error for login: %q", login)
		}
	}
}

func TestUserValidate_InvalidName(t *testing.T) {
	cases := []string{"",
		"a1234567890123456789012345678901234567890123456789012345678901234567890123456789012345678901234567890123456789012345678901234567891"}
	for _, name := range cases {
		user := &User{
			ID:           uuid.MustParse("123e4567-e89b-12d3-a456-426614174000"),
			Login:        "testuser",
			Name:         name,
			PasswordHash: []byte("hash"),
		}
		err := user.Validate()
		if err == nil {
			t.Errorf("Expected validation error for name: %q", name)
		}
	}
}

func TestUserValidate_EmptyPasswordHash(t *testing.T) {
	user := &User{
		ID:           uuid.MustParse("123e4567-e89b-12d3-a456-426614174000"),
		Login:        "testuser",
		Name:         "Test User",
		PasswordHash: []byte{},
	}

	err := user.Validate()
	if err == nil {
		t.Error("Expected validation error for empty password_hash")
	}
}
