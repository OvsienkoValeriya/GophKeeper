package models

import "time"

type ResourceType string

const (
	TypeCredentials ResourceType = "credentials"
	TypeText        ResourceType = "text"
	TypeBinary      ResourceType = "binary"
	TypeCard        ResourceType = "card"
)

type StorageType string

const (
	StoragePostgres StorageType = "postgres"
	StorageMinio    StorageType = "minio"
)

type Resource struct {
	ID        int64        `db:"id"`
	UserID    int64        `db:"user_id"`
	Name      string       `db:"name"`
	Type      ResourceType `db:"type"`
	Storage   StorageType  `db:"storage"`
	ObjectKey string       `db:"object_key"` // object key in MinIO if storage = minio
	Size      int64        `db:"size"`
	Metadata  []byte       `db:"metadata"` // encrypted metadata if storage = minio
	Data      []byte       `db:"data"`     // data if storage = postgres
	CreatedAt time.Time    `db:"created_at"`
	UpdatedAt time.Time    `db:"updated_at"`
}
