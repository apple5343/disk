package client

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"gitlab.crja72.ru/golang/2025/autumn/projects/go36/disk/services/auth-service/internal/model"
)

const (
	defaultClientTimeout = 5 * time.Second
	//nolint:gosec // Key name.
	InternalAPIKeyHeader = "X-Internal-API-Key"
)

type UserClient interface {
	CreateUser(ctx context.Context, user *model.UserCreateRequest) error
	GetUserByLogin(ctx context.Context, login string) (*model.InternalUser, error)
}

type HTTPUserClient struct {
	baseURL        string
	internalAPIKey string
	client         *http.Client
}

func NewHTTPUserClient(baseURL, internalAPIKey string) *HTTPUserClient {
	return &HTTPUserClient{
		baseURL:        baseURL,
		internalAPIKey: internalAPIKey,
		client: &http.Client{
			Timeout: defaultClientTimeout,
		},
	}
}

func (c *HTTPUserClient) CreateUser(ctx context.Context, user *model.UserCreateRequest) error {
	body, _ := json.Marshal(user)
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, c.baseURL+"/internal/users", bytes.NewReader(body))
	if err != nil {
		return err
	}
	req.Header.Set("Content-Type", "application/json")
	//nolint:canonicalheader // Key name.
	req.Header.Set(InternalAPIKeyHeader, c.internalAPIKey)

	resp, err := c.client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusConflict {
		return model.ErrUserExists
	}
	if resp.StatusCode != http.StatusCreated {
		return fmt.Errorf("user service returned %d", resp.StatusCode)
	}
	return nil
}

func (c *HTTPUserClient) GetUserByLogin(ctx context.Context, login string) (*model.InternalUser, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, c.baseURL+"/internal/users/by-login/"+login, nil)
	if err != nil {
		return nil, err
	}
	//nolint:canonicalheader // Key name.
	req.Header.Set(InternalAPIKeyHeader, c.internalAPIKey)

	resp, err := c.client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusNotFound {
		return nil, model.ErrUserNotFound
	}
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("user service returned %d", resp.StatusCode)
	}

	var user model.InternalUser
	if err = json.NewDecoder(resp.Body).Decode(&user); err != nil {
		return nil, err
	}
	return &user, nil
}
