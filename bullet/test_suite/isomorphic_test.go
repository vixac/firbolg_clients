package test_suite

import (
	"testing"

	"github.com/stretchr/testify/assert"
	ram "github.com/vixac/bullet/store/ram"
	"github.com/vixac/firbolg_clients/bullet/bullet_interface"
	local_bullet "github.com/vixac/firbolg_clients/bullet/local_bullet"
)

// the goal here is to test that both clients behave in the same way.
// The problem is I don't have a complete rest client setup, as each
// rest client test just sets up the 1 endpoint being tested.
// this can be added later.

func buildClients() []bullet_interface.BulletClientInterface {
	store := ram.NewRamStore()
	localClient := &local_bullet.LocalBullet{
		Store: store,
		AppId: 12,
	}
	var clients []bullet_interface.BulletClientInterface
	clients = append(clients, localClient)
	return clients
	//VX:TODO add rest client in here, and make this a map
}
func TestSomething(t *testing.T) {
	clients := buildClients()
	for _, c := range clients {
		err := c.TrackInsertOne(1, "testKey", int64(1234), nil, nil)
		err = c.TrackInsertOne(1, "testKey_2", int64(12345), nil, nil)
		assert.NoError(t, err)
		err = c.TrackInsertOne(1, "not_a_testKey3", int64(123456), nil, nil)
		assert.NoError(t, err)

		keys := []string{"testKey", "testKey_2"}
		//track get many
		bucket := bullet_interface.TrackGetKeys{
			BucketID: 1,
			Keys:     keys,
		}
		buckets := []bullet_interface.TrackGetKeys{bucket}
		req := bullet_interface.TrackGetManyRequest{
			Buckets: buckets,
		}
		res, err := c.TrackGetMany(req)
		assert.NoError(t, err)
		assert.NotNil(t, res)
		valuesInBucket, ok := res.Values[1]
		assert.True(t, ok)
		assert.Equal(t, len(valuesInBucket), 2)
		assert.Equal(t, valuesInBucket["testKey"].Value, int64(1234))
		assert.Equal(t, valuesInBucket["testKey_2"].Value, int64(12345))

		//trackgetmany by prefix

		prefixes := []string{"testKey_", "not_"}
		prefixReq := bullet_interface.TrackGetItemsbyManyPrefixesRequest{
			BucketID: 1,
			Prefixes: prefixes,
		}
		res, err = c.TrackGetByManyPrefixes(prefixReq)

		assert.NoError(t, err)
		assert.NotNil(t, res)
		valuesInBucket, ok = res.Values[1]
		assert.True(t, ok)
		assert.Equal(t, len(valuesInBucket), 2)
		assert.Equal(t, valuesInBucket["testKey_2"].Value, int64(12345))
		assert.Equal(t, valuesInBucket["not_a_testKey3"].Value, int64(123456))

	}

}
