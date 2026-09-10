package rest_bullet

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	bi "github.com/vixac/firbolg_clients/bullet/bullet_interface"
	"io"
	"net/http"
	"net/url"
)

var _ bi.WarehouseClientInterface = (*RestClient)(nil)

func (c *RestClient) warehouseJSON(ctx context.Context, method, path string, body, result any) error {
	data, err := json.Marshal(body)
	if err != nil {
		return err
	}
	req, err := http.NewRequestWithContext(ctx, method, c.BaseURL+path, bytes.NewReader(data))
	if err != nil {
		return err
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("X-App-Id", c.AppId)
	req.Header.Set("X-Tenancy-Id", c.TenancyId)
	if c.HTTPClient == nil {
		return fmt.Errorf("HTTPClient is nil")
	}
	resp, err := c.HTTPClient.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		var cause error
		switch resp.StatusCode {
		case http.StatusNotFound:
			cause = bi.ErrBlobNotFound
		case http.StatusConflict:
			cause = bi.ErrWarehousePutConflict
		case http.StatusNotImplemented:
			cause = bi.ErrWarehouseUnsupported

		}
		data, err := io.ReadAll(resp.Body)
		if err != nil {
			return err
		}
		if resp.StatusCode == http.StatusBadRequest {
			var response struct {
				Error string `json:"error"`
			}
			if json.Unmarshal(data, &response) == nil && response.Error == bi.ErrWarehouseInvalidPutID.Error() {
				cause = bi.ErrWarehouseInvalidPutID
			}
		}
		if cause != nil {
			return fmt.Errorf("warehouse HTTP %d: %s: %w", resp.StatusCode, data, cause)
		}
		return fmt.Errorf("warehouse HTTP %d: %s", resp.StatusCode, data)
	}
	return json.NewDecoder(resp.Body).Decode(result)
}
func (c *RestClient) WarehousePut(ctx context.Context, req bi.PutBlobRequest) (bi.Blob, error) {
	var b bi.Blob
	err := c.warehouseJSON(ctx, http.MethodPost, "/warehouse/blobs", req, &b)
	return b, err
}
func (c *RestClient) WarehouseGet(ctx context.Context, id bi.BlobID) (bi.Blob, error) {
	if id == "" {
		if err := ctx.Err(); err != nil {
			return bi.Blob{}, err
		}
		return bi.Blob{}, bi.ErrBlobNotFound
	}
	var b bi.Blob
	err := c.warehouseJSON(ctx, http.MethodGet, "/warehouse/blobs/"+url.PathEscape(string(id)), nil, &b)
	return b, err
}
func (c *RestClient) WarehouseGetMany(ctx context.Context, ids []bi.BlobID) (map[bi.BlobID]bi.Blob, error) {
	result := make(map[bi.BlobID]bi.Blob)
	err := c.warehouseJSON(ctx, http.MethodPost, "/warehouse/blobs/batch-get", struct {
		IDs []bi.BlobID `json:"ids"`
	}{ids}, &result)
	if err != nil {
		return nil, err
	}
	return result, nil
}
