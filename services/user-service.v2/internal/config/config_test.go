package config

import (
	"os"
	"testing"
)

func TestRead(t *testing.T) {
	oldUserAddr := os.Getenv("USER_SVC_ADDR")
	oldUserPort := os.Getenv("USER_SVC_PORT")
	oldDBURL := os.Getenv("USER_DB_DSN")
	oldAPIKey := os.Getenv("INTERNAL_API_KEY")

	t.Setenv("USER_SVC_ADDR", "test-user-service:9000")
	t.Setenv("USER_SVC_PORT", "9000")
	t.Setenv("USER_DB_DSN", "postgres://test:test@localhost:5433/test_db")
	t.Setenv("INTERNAL_API_KEY", "test-secret-key")

	defer func() {
		if oldUserAddr != "" {
			t.Setenv("USER_SVC_ADDR", oldUserAddr)
		} else {
			os.Unsetenv("USER_SVC_ADDR")
		}
		if oldUserPort != "" {
			t.Setenv("USER_SVC_PORT", oldUserPort)
		} else {
			os.Unsetenv("USER_SVC_PORT")
		}
		if oldDBURL != "" {
			t.Setenv("USER_DB_DSN", oldDBURL)
		} else {
			os.Unsetenv("USER_DB_DSN")
		}
		if oldAPIKey != "" {
			t.Setenv("INTERNAL_API_KEY", oldAPIKey)
		} else {
			os.Unsetenv("INTERNAL_API_KEY")
		}
	}()

	cfg, err := Read()
	if err != nil {
		t.Fatalf("Read() failed: %v", err)
	}

	if cfg.UserSvcAddr != "test-user-service:9000" {
		t.Errorf("UserSvcAddr = %q, want %q", cfg.UserSvcAddr, "test-user-service:9000")
	}
	if cfg.UserSvcPort != "9000" {
		t.Errorf("UserSvcPort = %q, want %q", cfg.UserSvcPort, "9000")
	}
	if cfg.DatabaseURL != "postgres://test:test@localhost:5433/test_db" {
		t.Errorf("DatabaseURL = %q, want %q", cfg.DatabaseURL, "postgres://test:test@localhost:5433/test_db")
	}
	if cfg.InternalAPIKey != "test-secret-key" {
		t.Errorf("InternalAPIKey = %q, want %q", cfg.InternalAPIKey, "test-secret-key")
	}
}

func TestRead_Defaults(t *testing.T) {
	oldUserAddr := os.Getenv("USER_SVC_ADDR")
	oldUserPort := os.Getenv("USER_SVC_PORT")
	oldDBURL := os.Getenv("USER_DB_DSN")
	oldAPIKey := os.Getenv("INTERNAL_API_KEY")

	os.Unsetenv("USER_SVC_ADDR")
	os.Unsetenv("USER_SVC_PORT")
	os.Unsetenv("USER_DB_DSN")
	os.Unsetenv("INTERNAL_API_KEY")

	defer func() {
		if oldUserAddr != "" {
			t.Setenv("USER_SVC_ADDR", oldUserAddr)
		}
		if oldUserPort != "" {
			t.Setenv("USER_SVC_PORT", oldUserPort)
		}
		if oldDBURL != "" {
			t.Setenv("USER_DB_DSN", oldDBURL)
		}
		if oldAPIKey != "" {
			t.Setenv("INTERNAL_API_KEY", oldAPIKey)
		}
	}()

	cfg, err := Read()
	if err != nil {
		t.Fatalf("Read() failed: %v", err)
	}

	if cfg.UserSvcAddr != "user-service:8002" {
		t.Errorf("UserSvcAddr default = %q, want %q", cfg.UserSvcAddr, "user-service:8002")
	}
	if cfg.UserSvcPort != "8002" {
		t.Errorf("UserSvcPort default = %q, want %q", cfg.UserSvcPort, "8002")
	}
	if cfg.DatabaseURL != "postgres://postgres:postgres@postgres:5432/user_db?sslmode=disable" {
		t.Errorf("DatabaseURL default = %q, want %q", cfg.DatabaseURL,
			"postgres://postgres:postgres@postgres:5432/user_db?sslmode=disable")
	}
	if cfg.InternalAPIKey != "secret-key-for-internal-communication" {
		t.Errorf("InternalAPIKey default = %q, want %q", cfg.InternalAPIKey, "secret-key-for-internal-communication")
	}
}
