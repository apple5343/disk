package repository

import (
	"errors"
	"testing"
)

func TestErrUserNotFound(t *testing.T) {
	if ErrUserNotFound == nil {
		t.Error("ErrUserNotFound must not be nil")
	}
	if ErrUserNotFound.Error() != "user not found" {
		t.Errorf("Unexpected error message: %s", ErrUserNotFound.Error())
	}
}

func TestNewPostgresRepository(t *testing.T) {
	repo := NewPostgresRepository(nil)
	if repo == nil {
		t.Error("NewPostgresRepository returned nil")
	}
}

func TestPostgresRepositoryImplementsUserRepository(_ *testing.T) {
	var _ UserRepository = &PostgresRepository{}
}

func TestErrorIs(t *testing.T) {
	err1 := ErrUserNotFound
	err2 := errors.New("user not found")

	if !errors.Is(err1, ErrUserNotFound) {
		t.Error("errors.Is failed for ErrUserNotFound")
	}

	if errors.Is(err2, ErrUserNotFound) {
		t.Error("errors.Is should not match different error instances")
	}
}
