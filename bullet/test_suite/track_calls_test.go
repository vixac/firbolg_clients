package test_suite

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	bi "github.com/vixac/firbolg_clients/bullet/bullet_interface"
)

func TestTrackGetAndPutMany(t *testing.T) {
	for _, pair := range buildClientPairs(t) {
		t.Run(pair.name, func(t *testing.T) {
			for name, client := range map[string]bi.BulletClientInterface{"local": pair.local, "rest": pair.rest} {
				t.Run(name, func(t *testing.T) {
					key := name + "/arbitrary ?&日本語"
					tag, metric := int64(42), 1.25
					large := int64(9007199254740993)
					require.NoError(t, client.TrackPutMany(bi.TrackPutManyRequest{Buckets: []bi.TrackPutItems{
						{BucketID: 1, Items: []bi.TrackKeyValueItem{{Key: key, Value: large, Tag: &tag, Metric: &metric}}},
						{BucketID: 2, Items: []bi.TrackKeyValueItem{{Key: key, Value: -large}}},
						{BucketID: 1, Items: []bi.TrackKeyValueItem{{Key: key + "zero", Value: 0}}},
					}}))
					for _, reader := range []bi.BulletClientInterface{pair.local, pair.rest} {
						value, err := reader.TrackGet(1, key)
						require.NoError(t, err)
						assert.Equal(t, large, value)
						value, err = reader.TrackGet(2, key)
						require.NoError(t, err)
						assert.Equal(t, -large, value)
						value, err = reader.TrackGet(1, key+"zero")
						require.NoError(t, err)
						assert.Zero(t, value)
						_, err = reader.TrackGet(1, key+"missing")
						require.Error(t, err)
						many, err := reader.TrackGetMany(bi.TrackGetManyRequest{Buckets: []bi.TrackGetKeys{{BucketID: 1, Keys: []string{key}}}})
						require.NoError(t, err)
						assert.Equal(t, bi.TrackValue{Value: large, Tag: &tag, Metric: &metric}, many.Values[1][key])
					}
					require.NoError(t, client.TrackPutMany(bi.TrackPutManyRequest{Buckets: []bi.TrackPutItems{{BucketID: 1, Items: []bi.TrackKeyValueItem{{Key: key, Value: 7}}}}}))
					value, err := client.TrackGet(1, key)
					require.NoError(t, err)
					assert.Equal(t, int64(7), value)
				})
			}
		})
	}
}

func TestTrackMutate(t *testing.T) {
	for _, pair := range buildClientPairs(t) {
		t.Run(pair.name, func(t *testing.T) {
			require.NoError(t, pair.local.TrackInsertOne(1, "delete", 1, nil, nil))
			tag, metric := int64(12), 3.5
			req := bi.TrackMutation{MutationID: "mutation-1", Puts: []bi.TrackRequest{
				{BucketID: 1, Key: "put", Value: 9007199254740993, Tag: &tag, Metric: &metric},
				{BucketID: 2, Key: "put", Value: -2},
			}, Deletes: []bi.TrackDeleteValue{{BucketID: 1, Key: "delete"}}}
			unsupported, err := pair.rest.TrackMutate(req)
			require.ErrorIs(t, err, bi.ErrTrackMutationUnsupported)
			assert.False(t, unsupported.Applied)
			value, err := pair.local.TrackGet(1, "delete")
			require.NoError(t, err)
			assert.Equal(t, int64(1), value)
			result, err := pair.local.TrackMutate(req)
			require.NoError(t, err)
			assert.True(t, result.Applied)
			_, err = pair.rest.TrackGet(1, "delete")
			require.Error(t, err)
			many, err := pair.rest.TrackGetMany(bi.TrackGetManyRequest{Buckets: []bi.TrackGetKeys{{BucketID: 1, Keys: []string{"put"}}, {BucketID: 2, Keys: []string{"put"}}}})
			require.NoError(t, err)
			assert.Equal(t, bi.TrackValue{Value: 9007199254740993, Tag: &tag, Metric: &metric}, many.Values[1]["put"])
			assert.Equal(t, int64(-2), many.Values[2]["put"].Value)
			require.NoError(t, pair.local.TrackInsertOne(1, "put", 99, nil, nil))
			require.NoError(t, pair.local.TrackInsertOne(1, "delete", 99, nil, nil))
			result, err = pair.local.TrackMutate(req)
			require.NoError(t, err)
			assert.False(t, result.Applied)
			for _, key := range []string{"put", "delete"} {
				value, err := pair.local.TrackGet(1, key)
				require.NoError(t, err)
				assert.Equal(t, int64(99), value)
			}
		})
	}
}
