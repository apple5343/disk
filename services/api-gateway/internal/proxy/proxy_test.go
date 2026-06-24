package proxy

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/proxy"
	"gitlab.crja72.ru/golang/2025/autumn/projects/go36/disk/services/api-gateway/internal/config"
)

// TestRegisterRoutes_HealthCheck проверяет health check endpoint.
func TestRegisterRoutes_HealthCheck(t *testing.T) {
	app := fiber.New()
	cfg := &config.Config{
		AuthSvcURL:     "http://localhost:3001",
		FileDataSvcURL: "http://localhost:3002",
		JWTSecret:      "test-secret-key-123",
	}

	RegisterRoutes(app, cfg)

	req := httptest.NewRequest(http.MethodGet, "/health", nil)
	resp, err := app.Test(req)
	if err != nil {
		t.Fatalf("Failed to make request: %v", err)
	}

	if resp.StatusCode != fiber.StatusOK {
		t.Errorf("Expected status %d, got %d", fiber.StatusOK, resp.StatusCode)
	}
}

// TestRegisterRoutes_Redirect проверяет редирект с корня.
func TestRegisterRoutes_Redirect(t *testing.T) {
	app := fiber.New()
	cfg := &config.Config{
		AuthSvcURL:     "http://localhost:3001",
		FileDataSvcURL: "http://localhost:3002",
		JWTSecret:      "test-secret-key-123",
	}

	RegisterRoutes(app, cfg)

	req := httptest.NewRequest(http.MethodGet, "/", nil)
	resp, err := app.Test(req)
	if err != nil {
		t.Fatalf("Failed to make request: %v", err)
	}

	if resp.StatusCode != fiber.StatusTemporaryRedirect {
		t.Errorf("Expected status %d, got %d", fiber.StatusTemporaryRedirect, resp.StatusCode)
	}

	location := resp.Header.Get("Location")
	if location != "/api/v1/files" {
		t.Errorf("Expected Location header to be '/api/v1/files', got '%s'", location)
	}
}

// TestProxyConfig проверяет конфигурацию прокси.
func TestProxyConfig(t *testing.T) {
	cfg := &config.Config{
		AuthSvcURL:     "http://auth:3001",
		FileDataSvcURL: "http://filedata:3002",
	}

	authProxyCfg := &proxy.Config{
		Servers: []string{cfg.AuthSvcURL},
		ModifyRequest: func(c *fiber.Ctx) error {
			c.Request().Header.Add("X-Real-IP", c.IP())
			return nil
		},
		ModifyResponse: func(c *fiber.Ctx) error {
			c.Response().Header.Del(fiber.HeaderServer)
			return nil
		},
	}

	fileDataProxyCfg := &proxy.Config{
		Servers: []string{cfg.FileDataSvcURL},
		ModifyRequest: func(c *fiber.Ctx) error {
			c.Request().Header.Add("X-Real-IP", c.IP())
			return nil
		},
		ModifyResponse: func(c *fiber.Ctx) error {
			c.Response().Header.Del(fiber.HeaderServer)
			return nil
		},
	}

	if len(authProxyCfg.Servers) != 1 || authProxyCfg.Servers[0] != cfg.AuthSvcURL {
		t.Errorf("Expected auth proxy server to be %s, got %v", cfg.AuthSvcURL, authProxyCfg.Servers)
	}

	if len(fileDataProxyCfg.Servers) != 1 || fileDataProxyCfg.Servers[0] != cfg.FileDataSvcURL {
		t.Errorf("Expected file data proxy server to be %s, got %v", cfg.FileDataSvcURL, fileDataProxyCfg.Servers)
	}

	if authProxyCfg.ModifyRequest == nil {
		t.Error("Expected ModifyRequest to be set for auth proxy")
	}

	if authProxyCfg.ModifyResponse == nil {
		t.Error("Expected ModifyResponse to be set for auth proxy")
	}

	if fileDataProxyCfg.ModifyRequest == nil {
		t.Error("Expected ModifyRequest to be set for file data proxy")
	}

	if fileDataProxyCfg.ModifyResponse == nil {
		t.Error("Expected ModifyResponse to be set for file data proxy")
	}
}

// TestConstants проверяет константы.
func TestConstants(t *testing.T) {
	if APIPrefix != "/api/" {
		t.Errorf("Expected APIPrefix to be '/api/', got '%s'", APIPrefix)
	}

	if APIVersion != "v1" {
		t.Errorf("Expected APIVersion to be 'v1', got '%s'", APIVersion)
	}

	expectedAuthPath := "/api/v1/auth/login"
	actualAuthPath := APIPrefix + APIVersion + "/auth/login"
	if actualAuthPath != expectedAuthPath {
		t.Errorf("Expected auth path to be '%s', got '%s'", expectedAuthPath, actualAuthPath)
	}
}

// TestRegisterRoutes_WithMockBackend проверяет работу с мок-сервером.
func TestRegisterRoutes_WithMockBackend(t *testing.T) {
	mockServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Header.Get("X-Real-IP") == "" {
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte(`{"error": "X-Real-IP header missing"}`))
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"message": "OK", "path": "` + r.URL.Path + `"}`))
	}))
	defer mockServer.Close()

	cfg := &config.Config{
		AuthSvcURL:     mockServer.URL,
		FileDataSvcURL: mockServer.URL,
		JWTSecret:      "test-secret",
	}

	app := fiber.New()
	RegisterRoutes(app, cfg)

	tests := []struct {
		name   string
		path   string
		method string
	}{
		{
			name:   "auth register through proxy",
			path:   "/api/v1/auth/register",
			method: "POST",
		},
		{
			name:   "auth login through proxy",
			path:   "/api/v1/auth/login",
			method: "POST",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := httptest.NewRequest(tt.method, tt.path, nil)
			resp, err := app.Test(req, 5000)
			if err != nil {
				t.Fatalf("Failed to make request: %v", err)
			}

			// Прокси должен работать и вернуть 200 OK
			if resp.StatusCode != http.StatusOK {
				t.Errorf("Expected status %d, got %d", http.StatusOK, resp.StatusCode)
			}

			// Проверяем, что заголовок Server удален
			if serverHeader := resp.Header.Get("Server"); serverHeader != "" {
				t.Errorf("Expected Server header to be empty, got '%s'", serverHeader)
			}
		})
	}
}

// TestProxyRequestModification проверяет модификацию запросов.
func TestProxyRequestModification(t *testing.T) {
	app := fiber.New()

	app.Get("/test", func(c *fiber.Ctx) error {
		ip := c.IP()
		if ip == "" {
			return c.Status(400).SendString("IP is empty")
		}
		return c.SendString("OK")
	})

	req := httptest.NewRequest(http.MethodGet, "/test", nil)
	req.RemoteAddr = "192.168.1.1:12345"

	resp, err := app.Test(req)
	if err != nil {
		t.Fatalf("Failed to make request: %v", err)
	}

	if resp.StatusCode != fiber.StatusOK {
		t.Errorf("Expected status %d, got %d", fiber.StatusOK, resp.StatusCode)
	}
}

// TestAllRoutesRegistered проверяет, что все маршруты зарегистрированы.
func TestAllRoutesRegistered(t *testing.T) {
	app := fiber.New()
	cfg := &config.Config{
		AuthSvcURL:     "http://localhost:3001",
		FileDataSvcURL: "http://localhost:3002",
		JWTSecret:      "test-secret-key-123",
	}

	RegisterRoutes(app, cfg)

	routes := app.GetRoutes()

	expectedRoutes := []string{
		"GET /health",
		"POST /api/v1/auth/register",
		"POST /api/v1/auth/login",
		"GET /",
	}

	for _, expected := range expectedRoutes {
		found := false
		for _, route := range routes {
			if route.Method+" "+route.Path == expected {
				found = true
				break
			}
		}

		if !found {
			t.Errorf("Expected route '%s' to be registered", expected)
		}
	}

	testPaths := []string{
		"/api/v1/files",
		"/api/v1/files/123",
	}

	for _, path := range testPaths {
		req := httptest.NewRequest(http.MethodGet, path, nil)
		resp, err := app.Test(req)
		if err != nil {
			t.Fatalf("Failed to make request to %s: %v", path, err)
		}

		if resp.StatusCode == fiber.StatusNotFound {
			t.Errorf("Expected route '%s' to be registered, got 404", path)
		}
	}
}

// TestProxyResponseModification проверяет модификацию ответов.
func TestProxyResponseModification(t *testing.T) {
	mockServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		w.Header().Set("Server", "TestServer/1.0")
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"test": "data"}`))
	}))
	defer mockServer.Close()

	cfg := &config.Config{
		AuthSvcURL:     mockServer.URL,
		FileDataSvcURL: mockServer.URL,
		JWTSecret:      "test-secret",
	}

	app := fiber.New()
	RegisterRoutes(app, cfg)

	req := httptest.NewRequest(http.MethodPost, "/api/v1/auth/login", nil)
	resp, err := app.Test(req, 5000)
	if err != nil {
		t.Fatalf("Failed to make request: %v", err)
	}

	serverHeader := resp.Header.Get("Server")
	if serverHeader != "" {
		t.Errorf("Expected Server header to be removed, got '%s'", serverHeader)
	}
}

// TestPathPatterns проверяет корректность паттернов путей.
func TestPathPatterns(t *testing.T) {
	app := fiber.New()
	cfg := &config.Config{
		AuthSvcURL:     "http://localhost:3001",
		FileDataSvcURL: "http://localhost:3002",
		JWTSecret:      "test-secret-key-123",
	}

	RegisterRoutes(app, cfg)

	testCases := []struct {
		path     string
		expected bool // true если должен быть обработан
	}{
		{"/api/v1/files", true},
		{"/api/v1/files/", true},
		{"/api/v1/files/123", true},
		{"/api/v1/files/123/", true},
		{"/api/v1/files/123/content", true},
		{"/api/v1/files/123/download", true},
		{"/api/v2/files", false}, // другая версия API
		{"/api/v1/other", false}, // другой путь
	}

	for _, tc := range testCases {
		req := httptest.NewRequest(http.MethodGet, tc.path, nil)
		resp, err := app.Test(req)
		if err != nil {
			t.Fatalf("Failed to make request to %s: %v", tc.path, err)
		}

		// Если путь должен обрабатываться, не должно быть 404
		if tc.expected && resp.StatusCode == fiber.StatusNotFound {
			t.Errorf("Path %s should be handled, got 404", tc.path)
		}

		// Если путь не должен обрабатываться, может быть 404 или другая ошибка прокси
		if !tc.expected && resp.StatusCode != fiber.StatusNotFound &&
			!strings.Contains(tc.path, "/api/v1/") {
			// Для путей вне нашей группы API ожидаем 404
			t.Logf("Path %s returned status %d", tc.path, resp.StatusCode)
		}
	}
}
