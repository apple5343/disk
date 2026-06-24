package sqlx

import (
	"context"
	"data/internal/models"
	"data/internal/repository"
	repoModels "data/internal/repository/file/models"
	"database/sql"
	"errors"
	"strconv"
	"strings"
)

func (r *fileRepository) GetFileByID(ctx context.Context, id string) (*models.FileMetadata, error) {
	var repoFile repoModels.FileMetadata

	err := r.db.GetContext(
		ctx,
		&repoFile,
		"SELECT id, user_id, storage_path, file_name, folder_id, bucket, full_path, size, content_type, tags, created_at, updated_at, status FROM files WHERE id = $1",
		id,
	)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, repository.ErrNotFound
		}
		return nil, err
	}

	result, err := repoModels.FileMetadataFromRepo(&repoFile)
	if err != nil {
		return nil, err
	}
	return result, nil
}

func (r *fileRepository) GetFileByStoragePath(ctx context.Context, storagePath string) (*models.FileMetadata, error) {
	var repoFile repoModels.FileMetadata

	err := r.db.GetContext(
		ctx,
		&repoFile,
		"SELECT id, user_id, storage_path, file_name, folder_id, bucket, full_path, size, content_type, tags, created_at, updated_at, status FROM files WHERE storage_path = $1",
		storagePath,
	)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, repository.ErrNotFound
		}
		return nil, err
	}

	result, err := repoModels.FileMetadataFromRepo(&repoFile)
	if err != nil {
		return nil, err
	}
	return result, nil
}

func (r *fileRepository) GetFileByPath(ctx context.Context, path string, userID string) (*models.FileMetadata, error) {
	var repoFile repoModels.FileMetadata

	err := r.db.GetContext(
		ctx,
		&repoFile,
		"SELECT id, user_id, storage_path, file_name, folder_id, bucket, full_path, size, content_type, tags, created_at, updated_at, status FROM files WHERE full_path = $1 AND user_id = $2",
		path,
		userID,
	)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, repository.ErrNotFound
		}
		return nil, err
	}

	result, err := repoModels.FileMetadataFromRepo(&repoFile)
	if err != nil {
		return nil, err
	}
	return result, nil
}

func (r *fileRepository) GetFilesByFolderID(ctx context.Context, folderID string) ([]*models.FileMetadata, error) {
	var repoFiles []repoModels.FileMetadata

	err := r.db.SelectContext(
		ctx,
		&repoFiles,
		"SELECT id, user_id, storage_path, file_name, folder_id, bucket, full_path, size, content_type, tags, created_at, updated_at, status FROM files WHERE folder_id = $1",
		folderID,
	)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return []*models.FileMetadata{}, nil
		}
		return nil, err
	}
	files := make([]*models.FileMetadata, len(repoFiles))
	for i, repoFile := range repoFiles {
		files[i], err = repoModels.FileMetadataFromRepo(&repoFile)
		if err != nil {
			return nil, err
		}
	}
	return files, nil
}

func (r *fileRepository) SearchFiles(
	ctx context.Context,
	userID string,
	req *models.SearchRequest,
) ([]*models.FileMetadata, error) {
	query := "SELECT id, user_id, storage_path, file_name, folder_id, bucket, full_path, size, content_type, tags, created_at, updated_at, status FROM files WHERE user_id = $1"
	params := []interface{}{userID}
	paramsLen := len(params)
	if req.FolderID != "" {
		query += " AND folder_id = $" + strconv.Itoa(paramsLen+1)
		params = append(params, req.FolderID)
		paramsLen++
	}
	if req.Path != "" {
		query += " AND full_path ILIKE $" + strconv.Itoa(paramsLen+1)
		params = append(params, strings.ReplaceAll(req.Path, "*", "%"))
		paramsLen++
	}
	if req.FileName != "" {
		query += " AND file_name ILIKE $" + strconv.Itoa(paramsLen+1)
		params = append(params, strings.ReplaceAll(req.FileName, "*", "%"))
	}
	var repoFiles []repoModels.FileMetadata
	err := r.db.SelectContext(ctx, &repoFiles, query, params...)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return []*models.FileMetadata{}, nil
		}
		return nil, err
	}
	files := make([]*models.FileMetadata, len(repoFiles))
	for i, repoFile := range repoFiles {
		files[i], err = repoModels.FileMetadataFromRepo(&repoFile)
		if err != nil {
			return nil, err
		}
	}
	return files, nil
}
