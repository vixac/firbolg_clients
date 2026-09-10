package bullet_interface

import (
	"context"
	"errors"
	"time"
)

type BlobID string
type PutID string

// PutBlobRequest identifies an immutable write within a tenancy space.
// Checksum is opaque metadata; it is not computed or validated yet.
type PutBlobRequest struct {
	PutID       PutID  `json:"put_id"`
	ContentType string `json:"content_type"`
	Value       []byte `json:"value"`
	Checksum    string `json:"checksum"`
}

type Blob struct {
	ID          BlobID    `json:"id"`
	PutID       PutID     `json:"put_id"`
	ContentType string    `json:"content_type"`
	Value       []byte    `json:"value"`
	Checksum    string    `json:"checksum"`
	CreatedAt   time.Time `json:"created_at"`
}

// WarehouseClientInterface uses the space configured on the client.
type WarehouseClientInterface interface {
	WarehousePut(context.Context, PutBlobRequest) (Blob, error)
	WarehouseGet(context.Context, BlobID) (Blob, error)
	WarehouseGetMany(context.Context, []BlobID) (map[BlobID]Blob, error)
}

var (
	ErrWarehouseUnsupported  = errors.New("warehouse is not supported by this store")
	ErrWarehouseInvalidPutID = errors.New("warehouse put ID must not be empty")
	ErrWarehousePutConflict  = errors.New("warehouse put ID already exists with different content or metadata")
	ErrBlobNotFound          = errors.New("blob not found")
)
