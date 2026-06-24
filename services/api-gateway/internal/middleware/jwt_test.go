package middleware

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
)

// TestJWTProtected_ValidToken проверяет обработку валидного jwt токена.
func TestJWTProtected_ValidToken(t *testing.T) {
	app := fiber.New()
	secretKey := "test-secret-key-12345"

	app.Get("/protected", JWTProtected(secretKey), func(c *fiber.Ctx) error {
		return c.JSON(fiber.Map{
			"message": "protected data",
			"success": true,
		})
	})

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"sub": "user123",
		"exp": time.Now().Add(time.Hour).Unix(),
		"iat": time.Now().Unix(),
	})

	tokenString, err := token.SignedString([]byte(secretKey))
	if err != nil {
		t.Fatalf("Failed to create JWT token: %v", err)
	}

	req := httptest.NewRequest(http.MethodGet, "/protected", nil)
	req.Header.Set("Authorization", "Bearer "+tokenString)

	resp, err := app.Test(req)
	if err != nil {
		t.Fatalf("Failed to make request: %v", err)
	}

	if resp.StatusCode != fiber.StatusOK {
		t.Errorf("Expected status %d for valid token, got %d", fiber.StatusOK, resp.StatusCode)
	}
}

// TestJWTProtected_MissingToken проверяет обработку запроса без токена.
func TestJWTProtected_MissingToken(t *testing.T) {
	app := fiber.New()
	secretKey := "test-secret-key"

	app.Get("/protected", JWTProtected(secretKey), func(c *fiber.Ctx) error {
		return c.SendString("should not reach here")
	})

	req := httptest.NewRequest(http.MethodGet, "/protected", nil)
	resp, err := app.Test(req)
	if err != nil {
		t.Fatalf("Failed to make request: %v", err)
	}

	if resp.StatusCode != fiber.StatusUnauthorized {
		t.Errorf("Expected status %d for missing token, got %d", fiber.StatusUnauthorized, resp.StatusCode)
	}
}

// TestJWTProtected_InvalidToken проверяет обработку невалидного токена.
func TestJWTProtected_InvalidToken(t *testing.T) {
	app := fiber.New()
	secretKey := "test-secret-key"

	app.Get("/protected", JWTProtected(secretKey), func(c *fiber.Ctx) error {
		return c.SendString("should not reach here")
	})

	testCases := []struct {
		name        string
		authHeader  string
		description string
	}{
		{
			name:        "malformed token",
			authHeader:  "Bearer invalid.token.here",
			description: "токен неправильного формата",
		},
		{
			name:        "wrong signing key",
			authHeader:  "Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJzdWIiOiIxMjM0NTY3ODkwIiwibmFtZSI6IkpvaG4gRG9lIiwiaWF0IjoxNTE2MjM5MDIyfQ.SflKxwRJSMeKKF2QT4fwpMeJf36POk6yJV_adQssw5c",
			description: "токен подписан другим ключом",
		},
		{
			name:        "no bearer prefix",
			authHeader:  "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJzdWIiOiIxMjM0NTY3ODkwIiwibmFtZSI6IkpvaG4gRG9lIiwiaWF0IjoxNTE2MjM5MDIyfQ.SflKxwRJSMeKKF2QT4fwpMeJf36POk6yJV_adQssw5c",
			description: "отсутствует префикс Bearer",
		},
		{
			name:        "empty token",
			authHeader:  "Bearer ",
			description: "пустой токен",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			req := httptest.NewRequest(http.MethodGet, "/protected", nil)
			req.Header.Set("Authorization", tc.authHeader)

			resp, err := app.Test(req)
			if err != nil {
				t.Fatalf("Failed to make request: %v", err)
			}

			if resp.StatusCode != fiber.StatusUnauthorized {
				t.Errorf("Expected status %d for %s, got %d",
					fiber.StatusUnauthorized, tc.description, resp.StatusCode)
			}
		})
	}
}

// TestJWTProtected_ExpiredToken проверяет обработку просроченного токена.
func TestJWTProtected_ExpiredToken(t *testing.T) {
	app := fiber.New()
	secretKey := "test-secret-key"

	app.Get("/protected", JWTProtected(secretKey), func(c *fiber.Ctx) error {
		return c.SendString("should not reach here")
	})

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"sub": "user123",
		"exp": time.Now().Add(-time.Hour).Unix(), // час назад
		"iat": time.Now().Add(-2 * time.Hour).Unix(),
	})

	tokenString, err := token.SignedString([]byte(secretKey))
	if err != nil {
		t.Fatalf("Failed to create expired JWT token: %v", err)
	}

	req := httptest.NewRequest(http.MethodGet, "/protected", nil)
	req.Header.Set("Authorization", "Bearer "+tokenString)

	resp, err := app.Test(req)
	if err != nil {
		t.Fatalf("Failed to make request: %v", err)
	}

	if resp.StatusCode != fiber.StatusUnauthorized {
		t.Errorf("Expected status %d for expired token, got %d",
			fiber.StatusUnauthorized, resp.StatusCode)
	}
}

// TestJWTError проверяет обработчик ошибок jwt.
func TestJWTError(t *testing.T) {
	app := fiber.New()

	app.Get("/error", func(c *fiber.Ctx) error {
		return jwtError(c, nil)
	})

	req := httptest.NewRequest(http.MethodGet, "/error", nil)
	resp, err := app.Test(req)
	if err != nil {
		t.Fatalf("Failed to make request: %v", err)
	}

	if resp.StatusCode != fiber.StatusUnauthorized {
		t.Errorf("Expected status %d from jwtError, got %d",
			fiber.StatusUnauthorized, resp.StatusCode)
	}

	body := make([]byte, resp.ContentLength)
	resp.Body.Read(body)

	expectedBody := `{"error":true,"msg":"Expired or invalid token"}`
	if !strings.Contains(string(body), expectedBody) {
		t.Errorf("Expected response body to contain %s, got %s", expectedBody, string(body))
	}
}
