package file

import (
	"context"
	"encoding/json"
	adapters "storage/internal/adapter"
	"storage/internal/infrastructure/kafka"
	"storage/internal/models"
)

const (
	UploadingFileTopic = "uploading_files"

	statusUploading = "uploading"
	statusUploaded  = "uploaded"
	statusFailed    = "failed"
)

type adapter struct {
	producer kafka.Producer
}

func NewAdapter(producer kafka.Producer) adapters.FileAdapter {
	return &adapter{
		producer: producer,
	}
}

func (a *adapter) PushUploadingFile(ctx context.Context, metadata *models.FileMetadata) error {
	value, err := json.Marshal(MessageToJSON(metadata, statusUploading, ""))
	if err != nil {
		return err
	}
	return a.producer.Publish(ctx, UploadingFileTopic, []byte(metadata.ID), value)
}

func (a *adapter) PushUploadedFile(ctx context.Context, metadata *models.FileMetadata) error {
	value, err := json.Marshal(MessageToJSON(metadata, statusUploaded, ""))
	if err != nil {
		return err
	}
	return a.producer.Publish(ctx, UploadingFileTopic, []byte(metadata.ID), value)
}

func (a *adapter) PushFailedFile(ctx context.Context, metadata *models.FileMetadata, err error) error {
	value, err := json.Marshal(MessageToJSON(metadata, statusFailed, err.Error()))
	if err != nil {
		return err
	}
	return a.producer.Publish(ctx, UploadingFileTopic, []byte(metadata.ID), value)
}
