package httpclient

import (
	"net/http"
	"storage/internal/config"
	filedata "storage/internal/infrastructure/file-data"
)

type client struct {
	baseURL string
	c       *http.Client
}

func NewClient(cfg *config.FileDataClientConfig) filedata.Client {
	return &client{
		baseURL: cfg.BaseURL,
		c:       &http.Client{},
	}
}
