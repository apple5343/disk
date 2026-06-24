package config

import (
	"testing"
)

// TestRead_UsesDefaults проверка дефолтных значений.
func TestRead_UsesDefaults(t *testing.T) {
	cfg, err := Read()
	if err != nil {
		t.Fatalf("Read() error = %v", err)
	}

	if cfg.Port != "8000" {
		t.Errorf("Port default = %q, want %q", cfg.Port, "8000")
	}
	if cfg.AuthSvcURL != "http://auth-service:8001" {
		t.Errorf("AuthSvcURL default = %q, want %q", cfg.AuthSvcURL, "http://auth-service:8001")
	}
	if cfg.FileDataSvcURL != "http://file-data-service:8003" {
		t.Errorf("FileDataSvcURL default = %q, want %q", cfg.FileDataSvcURL, "http://file-data-service:8003")
	}
	if cfg.JWTSecret != "super-secret-jwt-key-for-dev-only" {
		t.Errorf("JWTSecret default = %q, want %q", cfg.JWTSecret, "super-secret-jwt-key-for-dev-only")
	}
}

// TestRead_OverridesWithEnv проверка чтения переменных окружения.
func TestRead_OverridesWithEnv(t *testing.T) {
	t.Setenv("GW_PORT", "9000")
	t.Setenv("AUTH_SVC_URL", "http://auth:9001")
	t.Setenv("FILE_DATA_SVC_URL", "http://file-data:9003")
	t.Setenv("JWT_SECRET", "test-secret")

	cfg, err := Read()
	if err != nil {
		t.Fatalf("Read() error = %v", err)
	}

	if cfg.Port != "9000" {
		t.Errorf("Port = %q, want %q", cfg.Port, "9000")
	}
	if cfg.AuthSvcURL != "http://auth:9001" {
		t.Errorf("AuthSvcURL = %q, want %q", cfg.AuthSvcURL, "http://auth:9001")
	}
	if cfg.FileDataSvcURL != "http://file-data:9003" {
		t.Errorf("FileDataSvcURL = %q, want %q", cfg.FileDataSvcURL, "http://file-data:9003")
	}
	if cfg.JWTSecret != "test-secret" {
		t.Errorf("JWTSecret = %q, want %q", cfg.JWTSecret, "test-secret")
	}
}
