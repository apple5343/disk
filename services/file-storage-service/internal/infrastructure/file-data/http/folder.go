package httpclient

import (
	"context"
	"encoding/json"
	"net/http"
	filedata "storage/internal/infrastructure/file-data"
	"storage/internal/models"
	ctxutil "storage/internal/utils/ctx"

	"github.com/apple5343/errorx"
)

const (
	folderEndpoint = "/folders"
)

func (c *client) GetFolderByID(ctx context.Context, id string) (*models.Folder, error) {
	request, err := http.NewRequestWithContext(ctx, http.MethodGet, c.baseURL+folderEndpoint+"/"+id, nil)
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
		if resp.StatusCode == http.StatusBadRequest || resp.StatusCode == http.StatusNotFound {
			return nil, filedata.ErrNotFound
		}
		var e Error
		err = json.NewDecoder(resp.Body).Decode(&e)
		if err != nil {
			return nil, err
		}
		return nil, errorx.NewError(e.Message, errorx.Internal)
	}

	folder := &Folder{}
	err = json.NewDecoder(resp.Body).Decode(folder)
	if err != nil {
		return nil, err
	}
	return FolderFromHTTP(folder), nil
}
