package config

import (
	"os"
	"testing"
)

func TestGetConfigDefaults(t *testing.T) {
	keys := []string{"AUTH_SVC_PORT", "USER_SVC_URL", "USER_SVC_PORT", "JWT_SECRET", "INTERNAL_API_KEY"}
	for _, k := range keys {
		_ = os.Unsetenv(k)
	}

	cfg, err := Read()
	if err != nil {
		t.Fatalf("GetConfig returned error: %v", err)
	}

	if cfg.AuthSvcPort != "8001" {
		t.Fatalf("expected AuthSvcPort default 8001, got %q", cfg.AuthSvcPort)
	}
	if cfg.UserSvcURL != "http://user-service:8002" {
		t.Fatalf("expected UserSvcURL default http://user-service:8002, got %q", cfg.UserSvcURL)
	}
	if cfg.UserSvcPort != "8002" {
		t.Fatalf("expected UserSvcPort default 8002, got %q", cfg.UserSvcPort)
	}
	if cfg.JWTSecret == "" {
		t.Fatalf("expected JWTSecret be not empty, got %q", cfg.JWTSecret)
	}
	if cfg.InternalAPIKey != "secret-key-for-internal-communication" {
		t.Fatalf("expected InternalAPIKey default 8001, got %q", cfg.InternalAPIKey)
	}
}

func TestGetConfigFromEnv(t *testing.T) {
	t.Setenv("AUTH_SVC_PORT", "12345")
	t.Setenv("USER_SVC_URL", "http://custom:9000")
	t.Setenv("USER_SVC_PORT", "9000")
	t.Setenv("JWT_SECRET", "s3cr3t")
	t.Setenv("INTERNAL_API_KEY", "yoYoYo")

	cfg, err := Read()
	if err != nil {
		t.Fatalf("GetConfig returned error: %v", err)
	}

	if cfg.AuthSvcPort != "12345" {
		t.Fatalf("expected AuthSvcPort 12345, got %q", cfg.AuthSvcPort)
	}
	if cfg.UserSvcURL != "http://custom:9000" {
		t.Fatalf("expected UserSvcURL http://custom:9000, got %q", cfg.UserSvcURL)
	}
	if cfg.UserSvcPort != "9000" {
		t.Fatalf("expected UserSvcPort 9000, got %q", cfg.UserSvcPort)
	}
	if cfg.JWTSecret != "s3cr3t" {
		t.Fatalf("expected JWTSecret s3cr3t, got %q", cfg.JWTSecret)
	}
	if cfg.InternalAPIKey != "yoYoYo" {
		t.Fatalf("expected InternalAPIKey 8080, got %q", cfg.InternalAPIKey)
	}
}
