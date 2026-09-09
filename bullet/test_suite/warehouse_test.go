package test_suite

import (
	"context"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/require"
	"github.com/vixac/bullet/api"
	"github.com/vixac/bullet/store/ram"
	si "github.com/vixac/bullet/store/store_interface"
	bi "github.com/vixac/firbolg_clients/bullet/bullet_interface"
	"github.com/vixac/firbolg_clients/bullet/local_bullet"
	"github.com/vixac/firbolg_clients/bullet/rest_bullet"
	"net/http/httptest"
	"testing"
)

func TestWarehouseClients(t *testing.T) {
	s := ram.NewRamStore()
	server := httptest.NewServer(api.SetupWarehouseRouter(s, "/warehouse", gin.New()))
	defer server.Close()
	local := &local_bullet.LocalBullet{Store: s, Space: si.TenancySpace{AppId: 1, TenancyId: 2}}
	rest := rest_bullet.NewRestClient(server.URL, rest_bullet.FirbolgClientenancySpace{AppId: 1, TenancyId: 2})
	ctx := context.Background()
	req := bi.PutBlobRequest{PutID: "shared", Value: []byte{0, 255, 1}, ContentType: "application/octet-stream", Checksum: "opaque"}
	original, err := local.WarehousePut(ctx, req)
	require.NoError(t, err)
	for name, client := range map[string]bi.WarehouseClientInterface{"local": local, "rest": rest} {
		t.Run(name, func(t *testing.T) {
			retry, err := client.WarehousePut(ctx, req)
			require.NoError(t, err)
			require.Equal(t, original, retry)
			got, err := client.WarehouseGet(ctx, original.ID)
			require.NoError(t, err)
			require.Equal(t, original, got)
			many, err := client.WarehouseGetMany(ctx, []bi.BlobID{original.ID, "missing"})
			require.NoError(t, err)
			require.Equal(t, map[bi.BlobID]bi.Blob{original.ID: original}, many)
			_, err = client.WarehouseGet(ctx, "missing")
			require.ErrorIs(t, err, bi.ErrBlobNotFound)
			changed := req
			changed.Checksum = "different"
			_, err = client.WarehousePut(ctx, changed)
			require.ErrorIs(t, err, bi.ErrWarehousePutConflict)
			_, err = client.WarehousePut(ctx, bi.PutBlobRequest{})
			require.ErrorIs(t, err, bi.ErrWarehouseInvalidPutID)
			cancelled, cancel := context.WithCancel(ctx)
			cancel()
			_, err = client.WarehouseGet(ctx, original.ID)
			require.NoError(t, err)
			_, err = client.WarehousePut(cancelled, req)
			require.ErrorIs(t, err, context.Canceled)
			_, err = client.WarehouseGet(cancelled, original.ID)
			require.ErrorIs(t, err, context.Canceled)
			_, err = client.WarehouseGetMany(cancelled, nil)
			require.ErrorIs(t, err, context.Canceled)
		})
	}
}
