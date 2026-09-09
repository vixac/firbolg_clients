package local_bullet

import (
	"context"
	"errors"
	si "github.com/vixac/bullet/store/store_interface"
	bi "github.com/vixac/firbolg_clients/bullet/bullet_interface"
)

var _ bi.WarehouseClientInterface = (*LocalBullet)(nil)

func warehouseError(err error) error {
	switch {
	case errors.Is(err, si.ErrWarehouseUnsupported):
		return bi.ErrWarehouseUnsupported
	case errors.Is(err, si.ErrWarehouseInvalidPutID):
		return bi.ErrWarehouseInvalidPutID
	case errors.Is(err, si.ErrWarehousePutConflict):
		return bi.ErrWarehousePutConflict
	case errors.Is(err, si.ErrBlobNotFound):
		return bi.ErrBlobNotFound
	default:
		return err
	}
}
func warehouseBlob(b si.Blob) bi.Blob {
	return bi.Blob{ID: bi.BlobID(b.ID), PutID: bi.PutID(b.PutID), ContentType: b.ContentType, Value: b.Value, Checksum: b.Checksum, CreatedAt: b.CreatedAt}
}
func (l *LocalBullet) WarehousePut(ctx context.Context, req bi.PutBlobRequest) (bi.Blob, error) {
	b, err := l.Store.WarehousePut(ctx, l.Space, si.PutBlobRequest{PutID: si.PutID(req.PutID), ContentType: req.ContentType, Value: req.Value, Checksum: req.Checksum})
	return warehouseBlob(b), warehouseError(err)
}
func (l *LocalBullet) WarehouseGet(ctx context.Context, id bi.BlobID) (bi.Blob, error) {
	b, err := l.Store.WarehouseGet(ctx, l.Space, si.BlobID(id))
	return warehouseBlob(b), warehouseError(err)
}
func (l *LocalBullet) WarehouseGetMany(ctx context.Context, ids []bi.BlobID) (map[bi.BlobID]bi.Blob, error) {
	keys := make([]si.BlobID, len(ids))
	for i, id := range ids {
		keys[i] = si.BlobID(id)
	}
	found, err := l.Store.WarehouseGetMany(ctx, l.Space, keys)
	if err != nil {
		return nil, warehouseError(err)
	}
	result := make(map[bi.BlobID]bi.Blob, len(found))
	for id, b := range found {
		result[bi.BlobID(id)] = warehouseBlob(b)
	}
	return result, nil
}
