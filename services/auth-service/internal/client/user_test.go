package client

import (
	"context"
	"testing"

	"gitlab.crja72.ru/golang/2025/autumn/projects/go36/disk/services/auth-service/internal/model"
)

func TestNewHTTPUserClient(t *testing.T) {
	client := NewHTTPUserClient("http://test", "test-key")

	if client.baseURL != "http://test" {
		t.Error("BaseURL not set correctly")
	}
	if client.internalAPIKey != "test-key" {
		t.Error("InternalAPIKey not set correctly")
	}
}

func TestHTTPUserClient_CreateUser(t *testing.T) {
	client := NewHTTPUserClient("http://invalid", "test-key")

	req := &model.UserCreateRequest{
		ID:           "123",
		Login:        "test",
		Name:         "Test",
		PasswordHash: []byte("hash"),
	}

	err := client.CreateUser(context.Background(), req)
	if err != nil {
		t.Logf("Expected error (connection failure): %v", err)
	}
}

func TestHTTPUserClient_GetUserByLogin(t *testing.T) {
	client := NewHTTPUserClient("http://invalid", "test-key")

	user, err := client.GetUserByLogin(context.Background(), "testuser")
	if err != nil {
		t.Logf("Expected error (connection failure): %v", err)
	}
	if user != nil {
		t.Error("Expected nil user")
	}
}

func TestHTTPUserClientImplementsUserClient(_ *testing.T) {
	var _ UserClient = &HTTPUserClient{}
}
