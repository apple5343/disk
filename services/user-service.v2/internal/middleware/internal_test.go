package middleware

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gofiber/fiber/v2"
	"gitlab.crja72.ru/golang/2025/autumn/projects/go36/disk/services/user-service.v2/internal/config"
)

func TestInternalAuth(t *testing.T) {
	cfg := &config.Config{
		InternalAPIKey: "test-secret-key",
	}

	app := fiber.New()
	app.Use(InternalAuth(cfg))
	app.Get("/test", func(c *fiber.Ctx) error {
		return c.SendString("OK")
	})

	// Bad key.
	req1 := httptest.NewRequest(http.MethodGet, "/test", nil)
	//nolint:canonicalheader // As is.
	req1.Header.Set(InternalAPIKeyHeader, "test-secret-key")
	resp1, _ := app.Test(req1)
	if resp1.StatusCode != http.StatusOK {
		t.Errorf("expected 200, got %d", resp1.StatusCode)
	}

	// Bad key v2.
	req2 := httptest.NewRequest(http.MethodGet, "/test", nil)
	//nolint:canonicalheader // As is.
	req2.Header.Set(InternalAPIKeyHeader, "wrong-key")
	resp2, _ := app.Test(req2)
	if resp2.StatusCode != http.StatusForbidden {
		t.Errorf("expected 403, got %d", resp2.StatusCode)
	}

	// No hdr.
	req3 := httptest.NewRequest(http.MethodGet, "/test", nil)
	resp3, _ := app.Test(req3)
	if resp3.StatusCode != http.StatusForbidden {
		t.Errorf("expected 403, got %d", resp3.StatusCode)
	}
}
