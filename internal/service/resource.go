package service

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"io"

	"github.com/OvsienkoValeriya/GophKeeper/internal/logger"
	"github.com/OvsienkoValeriya/GophKeeper/internal/models"
	"github.com/OvsienkoValeriya/GophKeeper/internal/repository"
	"github.com/OvsienkoValeriya/GophKeeper/internal/repository/storage"
	"github.com/google/uuid"
	"github.com/minio/minio-go/v7"
)

var (
	ErrAccessDenied     = errors.New("access denied")
	ErrResourceNotFound = errors.New("resource not found")
)

const maxPostgresSize = 1 << 20 // 1 МБ

type ResourceService struct {
	resourceRepo repository.ResourceRepository
	fileStorage  storage.Storage
}

func NewResourceService(resourceRepo repository.ResourceRepository, fileStorage storage.Storage) *ResourceService {
	return &ResourceService{
		resourceRepo: resourceRepo,
		fileStorage:  fileStorage,
	}
}

// Upload uploads a resource:
// - small data (< 1 MB) is saved in PostgreSQL
// - large data (>= 1 MB) is saved in MinIO, metadata is saved in PostgreSQL
func (s *ResourceService) Upload(ctx context.Context, userID int64, name string,
	resourceType models.ResourceType, data []byte) (*models.Resource, error) {

	resource := &models.Resource{
		UserID: userID,
		Name:   name,
		Type:   resourceType,
		Size:   int64(len(data)),
	}

	if len(data) < maxPostgresSize {
		resource.Storage = models.StoragePostgres
		resource.Data = data
	} else {
		resource.Storage = models.StorageMinio
		resource.ObjectKey = generateObjectKey(userID)

		if err := s.fileStorage.Upload(ctx, resource.ObjectKey, bytes.NewReader(data), resource.Size, minio.PutObjectOptions{}); err != nil {
			return nil, fmt.Errorf("failed to upload to file storage: %w", err)
		}
	}

	created, err := s.resourceRepo.Create(ctx, resource)
	if err != nil {
		// Rollback: if the database write failed, delete from MinIO
		if resource.Storage == models.StorageMinio {
			_ = s.fileStorage.Delete(ctx, resource.ObjectKey, minio.RemoveObjectOptions{})
			logger.Sugar.Error("failed to delete from file storage", "error", err)
		}
		return nil, fmt.Errorf("failed to save resource: %w", err)
	}

	return created, nil
}

func (s *ResourceService) Get(ctx context.Context, userID, resourceID int64) (*models.Resource, []byte, error) {
	resource, err := s.resourceRepo.GetByID(ctx, resourceID)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to get resource: %w", err)
	}

	if resource.UserID != userID {
		return nil, nil, ErrAccessDenied
	}

	var data []byte

	if resource.Storage == models.StoragePostgres {
		data = resource.Data
	} else {
		reader, err := s.fileStorage.Download(ctx, resource.ObjectKey, minio.GetObjectOptions{})
		if err != nil {
			return nil, nil, fmt.Errorf("failed to download from file storage: %w", err)
		}
		defer reader.Close()

		data, err = io.ReadAll(reader)
		if err != nil {
			return nil, nil, fmt.Errorf("failed to read data: %w", err)
		}
	}

	return resource, data, nil
}

func (s *ResourceService) GetByName(ctx context.Context, userID int64, name string) (*models.Resource, []byte, error) {
	resource, err := s.resourceRepo.GetByNameAndUserID(ctx, userID, name)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to get resource: %w", err)
	}

	var data []byte

	if resource.Storage == models.StoragePostgres {
		data = resource.Data
	} else {
		reader, err := s.fileStorage.Download(ctx, resource.ObjectKey, minio.GetObjectOptions{})
		if err != nil {
			return nil, nil, fmt.Errorf("failed to download from file storage: %w", err)
		}
		defer reader.Close()

		data, err = io.ReadAll(reader)
		if err != nil {
			return nil, nil, fmt.Errorf("failed to read data: %w", err)
		}
	}

	return resource, data, nil
}

func (s *ResourceService) GetAll(ctx context.Context, userID int64) ([]*models.Resource, error) {
	resources, err := s.resourceRepo.GetByUserID(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("failed to get resources: %w", err)
	}

	return resources, nil
}

func (s *ResourceService) Update(ctx context.Context, userID, resourceID int64, name string,
	resourceType models.ResourceType, data []byte) (*models.Resource, error) {

	existing, err := s.resourceRepo.GetByID(ctx, resourceID)
	if err != nil {
		return nil, fmt.Errorf("failed to get resource: %w", err)
	}

	if existing.UserID != userID {
		return nil, ErrAccessDenied
	}

	newSize := int64(len(data))
	newStorage := models.StoragePostgres
	if len(data) >= maxPostgresSize {
		newStorage = models.StorageMinio
	}

	oldStorage := existing.Storage
	oldObjectKey := existing.ObjectKey

	resource := &models.Resource{
		ID:     resourceID,
		UserID: userID,
		Name:   name,
		Type:   resourceType,
		Size:   newSize,
	}

	if newStorage == models.StoragePostgres {
		resource.Storage = models.StoragePostgres
		resource.Data = data
		resource.ObjectKey = ""
	} else {
		resource.Storage = models.StorageMinio

		if oldStorage == models.StorageMinio && oldObjectKey != "" {
			resource.ObjectKey = oldObjectKey
		} else {
			resource.ObjectKey = generateObjectKey(userID)
		}

		if err := s.fileStorage.Upload(ctx, resource.ObjectKey, bytes.NewReader(data), newSize, minio.PutObjectOptions{}); err != nil {
			return nil, fmt.Errorf("failed to upload to file storage: %w", err)
		}
	}

	if err := s.resourceRepo.Update(ctx, resource); err != nil {
		// Rollback: if the database write failed and we uploaded a new file
		if newStorage == models.StorageMinio && oldStorage != models.StorageMinio {
			_ = s.fileStorage.Delete(ctx, resource.ObjectKey, minio.RemoveObjectOptions{})
		}
		return nil, fmt.Errorf("failed to update resource: %w", err)
	}

	// If it was in MinIO before and now in PostgreSQL, delete the old file
	if oldStorage == models.StorageMinio && newStorage == models.StoragePostgres && oldObjectKey != "" {
		_ = s.fileStorage.Delete(ctx, oldObjectKey, minio.RemoveObjectOptions{})
	}

	return resource, nil
}

func (s *ResourceService) Delete(ctx context.Context, userID, resourceID int64) error {
	resource, err := s.resourceRepo.GetByID(ctx, resourceID)
	if err != nil {
		return fmt.Errorf("failed to get resource: %w", err)
	}

	if resource.UserID != userID {
		return ErrAccessDenied
	}

	if resource.Storage == models.StorageMinio && resource.ObjectKey != "" {
		if err := s.fileStorage.Delete(ctx, resource.ObjectKey, minio.RemoveObjectOptions{}); err != nil {
			return fmt.Errorf("failed to delete from file storage: %w", err)
		}
	}

	if err := s.resourceRepo.Delete(ctx, resourceID); err != nil {
		return fmt.Errorf("failed to delete resource: %w", err)
	}

	return nil
}

func generateObjectKey(userID int64) string {
	return fmt.Sprintf("users/%d/%s", userID, uuid.New().String())
}
