package local_bullet_test

import (
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/vixac/bullet/model"
	"github.com/vixac/bullet/store/ram"
	"github.com/vixac/firbolg_clients/bullet/local_bullet"
)

func TestTrackPayloadIsExposed(t *testing.T) {
	client := local_bullet.NewLocalClient(ram.NewRamStore(), model.TenancySpace{AppId: 12, TenancyId: 100})
	want := []byte("track payload")

	require.NoError(t, client.TrackPut(42, "payload-key", model.TrackValue{
		Value:   7,
		Payload: want,
	}))

	withoutPayload, err := client.TrackGet(42, "payload-key", model.TrackReadOptions{})
	require.NoError(t, err)
	require.Nil(t, withoutPayload.Payload)

	withPayload, err := client.TrackGet(42, "payload-key", model.TrackReadOptions{IncludePayload: true})
	require.NoError(t, err)
	require.Equal(t, want, withPayload.Payload)
}
