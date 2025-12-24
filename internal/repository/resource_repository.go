package repository

import (
	"context"

	"github.com/OvsienkoValeriya/GophKeeper/internal/models"
)

type ResourceRepository interface {
	Create(ctx context.Context, resource *models.Resource) (*models.Resource, error)

	GetByID(ctx context.Context, id int64) (*models.Resource, error)

	GetByUserID(ctx context.Context, userID int64) ([]*models.Resource, error)

	GetByNameAndUserID(ctx context.Context, userID int64, name string) (*models.Resource, error)

	Update(ctx context.Context, resource *models.Resource) error

	Delete(ctx context.Context, id int64) error
}
