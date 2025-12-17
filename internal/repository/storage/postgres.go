package storage

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"github.com/OvsienkoValeriya/GophKeeper/internal/models"
	"github.com/jmoiron/sqlx"
)

type PostgresResourceRepository struct {
	db *sqlx.DB
}

var (
	ErrResourceNotFound = errors.New("resource not found")
)

func NewPostgresResourceRepository(dsn string) (*PostgresResourceRepository, error) {
	db, err := sqlx.Connect("pgx", dsn)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to db: %w", err)
	}
	return &PostgresResourceRepository{db: db}, nil
}

func (r *PostgresResourceRepository) Create(ctx context.Context, resource *models.Resource) (*models.Resource, error) {
	query := `
        INSERT INTO resources (user_id, name, type, storage, object_key, size, metadata, data)
        VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
        RETURNING id, created_at, updated_at
    `

	err := r.db.QueryRowxContext(ctx, query,
		resource.UserID,
		resource.Name,
		resource.Type,
		resource.Storage,
		resource.ObjectKey,
		resource.Size,
		resource.Metadata,
		resource.Data,
	).Scan(&resource.ID, &resource.CreatedAt, &resource.UpdatedAt)

	if err != nil {
		return nil, fmt.Errorf("failed to create resource: %w", err)
	}

	return resource, nil
}

func (r *PostgresResourceRepository) GetByID(ctx context.Context, id int64) (*models.Resource, error) {
	query := `
        SELECT id, user_id, name, type, storage, object_key, size, metadata, data, created_at, updated_at
        FROM resources
        WHERE id = $1
    `

	var resource models.Resource
	err := r.db.GetContext(ctx, &resource, query, id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrResourceNotFound
		}
		return nil, fmt.Errorf("failed to get resource: %w", err)
	}

	return &resource, nil
}

func (r *PostgresResourceRepository) GetByUserID(ctx context.Context, userID int64) ([]*models.Resource, error) {
	query := `
		SELECT id, user_id, name, type, storage, object_key, size, metadata, data, created_at, updated_at
		FROM resources
		WHERE user_id = $1
		ORDER BY created_at DESC
	`

	var resources []*models.Resource
	err := r.db.SelectContext(ctx, &resources, query, userID)
	if err != nil {
		return nil, fmt.Errorf("failed to get resources: %w", err)
	}

	return resources, nil
}

func (r *PostgresResourceRepository) GetByNameAndUserID(ctx context.Context, userID int64, name string) (*models.Resource, error) {
	query := `
		SELECT id, user_id, name, type, storage, object_key, size, metadata, data, created_at, updated_at
		FROM resources
		WHERE user_id = $1 AND name = $2
	`

	var resource models.Resource
	err := r.db.GetContext(ctx, &resource, query, userID, name)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrResourceNotFound
		}
		return nil, fmt.Errorf("failed to get resource by name: %w", err)
	}
	return &resource, nil
}

func (r *PostgresResourceRepository) Update(ctx context.Context, resource *models.Resource) error {
	query := `
		UPDATE resources
		SET name = $1, type = $2, storage = $3, object_key = $4, size = $5, metadata = $6, data = $7
		WHERE id = $8
	`

	_, err := r.db.ExecContext(ctx, query, resource.Name, resource.Type, resource.Storage, resource.ObjectKey, resource.Size, resource.Metadata, resource.Data, resource.ID)
	if err != nil {
		return fmt.Errorf("failed to update resource: %w", err)
	}
	return nil
}

func (r *PostgresResourceRepository) Delete(ctx context.Context, id int64) error {
	query := `
		DELETE FROM resources
		WHERE id = $1
	`

	_, err := r.db.ExecContext(ctx, query, id)
	if err != nil {
		return fmt.Errorf("failed to delete resource: %w", err)
	}
	return nil
}
