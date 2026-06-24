package httpclient

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	filedata "storage/internal/infrastructure/file-data"
	"storage/internal/models"
	ctxutil "storage/internal/utils/ctx"

	"github.com/apple5343/errorx"
)

const (
	fileMetadataEndpoint = "/files"
)

func (c *client) SaveFileMetadata(ctx context.Context, metadata *models.FileMetadata) (*models.FileMetadata, error) {
	data, err := json.Marshal(FileMetadataToHTTP(metadata))

	if err != nil {
		return nil, err
	}
	request, err := http.NewRequestWithContext(
		ctx,
		http.MethodPost,
		c.baseURL+fileMetadataEndpoint,
		bytes.NewReader(data),
	)
	request.Header.Set("Content-Type", "application/json")
	request.Header.Set("Authorization", "Bearer "+ctxutil.TokenFromContext(ctx))
	if err != nil {
		return nil, err
	}
	resp, err := c.c.Do(request)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		var e Error
		err = json.NewDecoder(resp.Body).Decode(&e)
		if err != nil {
			return nil, err
		}
		return nil, errorx.NewError(e.Message, errorx.Internal)
	}

	fileMetadata := &FileMetadata{}
	err = json.NewDecoder(resp.Body).Decode(fileMetadata)
	if err != nil {
		return nil, err
	}

	return FileMetadataFromHTTP(fileMetadata), nil
}

func (c *client) GetFileMetadata(ctx context.Context, id string) (*models.FileMetadata, error) {
	request, err := http.NewRequestWithContext(ctx, http.MethodGet, c.baseURL+fileMetadataEndpoint+"/"+id, nil)
	if err != nil {
		return nil, err
	}
	request.Header.Set("Authorization", "Bearer "+ctxutil.TokenFromContext(ctx))

	resp, err := c.c.Do(request)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		if resp.StatusCode == http.StatusBadRequest {
			return nil, filedata.ErrNotFound
		}

		var e Error
		err = json.NewDecoder(resp.Body).Decode(&e)
		if err != nil {
			return nil, err
		}
		return nil, errorx.NewError(e.Message, errorx.Internal)
	}

	fileMetadata := &FileMetadata{}
	err = json.NewDecoder(resp.Body).Decode(fileMetadata)
	if err != nil {
		return nil, err
	}

	return FileMetadataFromHTTP(fileMetadata), nil
}

func (c *client) DeleteFileMetadata(ctx context.Context, id string) error {
	request, err := http.NewRequestWithContext(ctx, http.MethodDelete, c.baseURL+fileMetadataEndpoint+"/"+id, nil)
	if err != nil {
		return err
	}
	request.Header.Set("Authorization", "Bearer "+ctxutil.TokenFromContext(ctx))

	resp, err := c.c.Do(request)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusNoContent {
		if resp.StatusCode == http.StatusBadRequest {
			return filedata.ErrNotFound
		}
		var e Error
		err = json.NewDecoder(resp.Body).Decode(&e)
		if err != nil {
			return err
		}
		return errorx.NewError(e.Message, errorx.Internal)
	}

	return nil
}
