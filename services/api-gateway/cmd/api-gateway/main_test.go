package main

import (
	"context"
	"errors"
	"testing"
	"time"

	"gitlab.crja72.ru/golang/2025/autumn/projects/go36/disk/services/api-gateway/internal/config"
)

func TestApp_Run_SuccessfulStartAndShutdown(t *testing.T) {
	cfg := &config.Config{Port: "0"}
	app := NewApp(cfg)

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	done := make(chan error, 1)
	go func() {
		done <- app.run(ctx, 5*time.Second)
	}()

	time.Sleep(100 * time.Millisecond)

	cancel()

	select {
	case err := <-done:
		if err != nil {
			t.Errorf("run() failed on shutdown: %v", err)
		}
	case <-time.After(2 * time.Second):
		t.Fatal("run() did not shutdown in time")
	}
}

func mockReadConfig() (*config.Config, error) {
	return &config.Config{Port: "8081"}, nil
}

func mockNewApp(cfg *config.Config) *App {
	return NewApp(cfg)
}

func mockReadConfigError() (*config.Config, error) {
	return nil, errors.New("config read error")
}

// TestTestableMain_Success - Тест successful testableMain.
func TestTestableMain_Success(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	done := make(chan error, 1)
	go func() {
		done <- testableMain(ctx, 5*time.Second, mockReadConfig, mockNewApp)
	}()
	time.Sleep(100 * time.Millisecond)
	cancel()

	select {
	case err := <-done:
		if err != nil {
			t.Errorf("testableMain() failed: %v", err)
		}
	case <-time.After(2 * time.Second):
		t.Fatal("testableMain() did not exit")
	}
}

// TestTestableMain_ConfigError - тест ошибки в readConfig.
func TestTestableMain_ConfigError(t *testing.T) {
	ctx := context.Background()

	err := testableMain(ctx, 5*time.Second, mockReadConfigError, mockNewApp)
	if err == nil || err.Error() != "config read error" {
		t.Errorf("Expected config error, got: %v", err)
	}
}
