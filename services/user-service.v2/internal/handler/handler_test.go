package handler

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gofiber/fiber/v2"
	"gitlab.crja72.ru/golang/2025/autumn/projects/go36/disk/services/user-service.v2/internal/config"
	"gitlab.crja72.ru/golang/2025/autumn/projects/go36/disk/services/user-service.v2/internal/service"
)

func TestRegisterHandlers(t *testing.T) {
	app := fiber.New()

	cfg := &config.Config{
		InternalAPIKey: "test-key",
	}

	userService := service.NewUserService(nil)

	RegisterHandlers(app, cfg, userService)

	testRoutes := []struct {
		method string
		path   string
	}{
		{"GET", "/health"},
		{"GET", "/api/v1/users/123"},
		{"POST", "/internal/users"},
		{"GET", "/internal/users/by-login/test"},
	}

	for _, tr := range testRoutes {
		req := httptest.NewRequest(tr.method, tr.path, nil)
		if tr.method == http.MethodPost {
			req.Header.Set("Content-Type", "application/json")
		}
		resp, err := app.Test(req)
		if err != nil {
			t.Fatalf("Test request failed: %v", err)
		}
		if resp.StatusCode == http.StatusNotFound {
			t.Errorf("Route %s %s returned 404 (not registered)", tr.method, tr.path)
		}
	}
}
