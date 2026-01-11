package test_suite

import (
	"fmt"
	"log"
	"path/filepath"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	ram "github.com/vixac/bullet/store/ram"
	sqlite_store "github.com/vixac/bullet/store/sqlite"
	"github.com/vixac/bullet/store/store_interface"
	"github.com/vixac/firbolg_clients/bullet/bullet_interface"
	local_bullet "github.com/vixac/firbolg_clients/bullet/local_bullet"
)

// the goal here is to test that both clients behave in the same way.
// The problem is I don't have a complete rest client setup, as each
// rest client test just sets up the 1 endpoint being tested.
// this can be added later.

func buildClients() []bullet_interface.BulletClientInterface {
	store := ram.NewRamStore()
	space := store_interface.TenancySpace{
		AppId:     12,
		TenancyId: 100,
	}
	localClient := &local_bullet.LocalBullet{
		Store: store,
		Space: space,
	}
	var clients []bullet_interface.BulletClientInterface
	clients = append(clients, localClient)

	sqlLiteStore, err := sqlite_store.NewSQLiteStore("test-sqlite")
	if err != nil {
		log.Fatal(err)
	}
	localSqlClient := &local_bullet.LocalBullet{
		Store: sqlLiteStore,
		Space: space,
	}
	clients = append(clients, localSqlClient)

	return clients
	//VX:TODO add rest client in here, and make this a map
}
func TestTrack(t *testing.T) {
	clients := buildClients()
	for _, c := range clients {
		err := c.TrackInsertOne(1, "testKey", int64(1234), nil, nil)
		assert.NoError(t, err)
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

func TestDepot(t *testing.T) {
	clients := buildClients()
	for _, c := range clients {
		err := c.DepotInsertOne(bullet_interface.DepotRequest{
			Key:   1,
			Value: "value1",
		})
		assert.NoError(t, err)
		//insert many, and overwrite key 1 too.
		many := []bullet_interface.DepotRequest{
			{
				Key:   1,
				Value: "new_value1",
			},
			{
				Key:   2,
				Value: "value2",
			},
			{
				Key:   3,
				Value: "value3",
			},
			{
				Key:   4,
				Value: "value4",
			},
		}

		err = c.DepotUpsertMany(many)
		assert.NoError(t, err)
		keys := []int64{1, 3, 10}
		manyReq := bullet_interface.DepotGetManyRequest{
			Keys: keys,
		}
		res, err := c.DepotGetMany(manyReq)
		assert.NoError(t, err)
		assert.Equal(t, len(res.Values), 2)
		assert.Equal(t, len(res.Missing), 1)
		assert.Equal(t, res.Values[1], "new_value1")
		assert.Equal(t, res.Values[3], "value3")
		assert.Equal(t, res.Missing[0], int64(10))
	}

}

func buildClientsForGrove(t *testing.T) []bullet_interface.BulletClientInterface {
	store := ram.NewRamStore()
	space := store_interface.TenancySpace{
		AppId:     12,
		TenancyId: 100,
	}
	localClient := &local_bullet.LocalBullet{
		Store: store,
		Space: space,
	}
	var clients []bullet_interface.BulletClientInterface
	clients = append(clients, localClient)

	// Use a unique temporary database file for this test run
	dbPath := filepath.Join(t.TempDir(), fmt.Sprintf("test-sqlite-%d", time.Now().UnixNano()))
	sqlLiteStore, err := sqlite_store.NewSQLiteStore(dbPath)
	if err != nil {
		log.Fatal(err)
	}
	localSqlClient := &local_bullet.LocalBullet{
		Store: sqlLiteStore,
		Space: space,
	}
	clients = append(clients, localSqlClient)

	return clients
}

func TestGrove(t *testing.T) {
	clients := buildClientsForGrove(t)
	for _, c := range clients {
		// Create a simple tree structure:
		//       root
		//      /    \
		//   child1  child2
		//    /
		// grandchild1

		// Create root node
		err := c.GroveCreateNode(bullet_interface.GroveCreateNodeRequest{
			NodeID:   "root",
			Parent:   nil,
			Position: nil,
			Metadata: nil,
		})
		assert.NoError(t, err)

		// Verify root exists
		existsRes, err := c.GroveExists(bullet_interface.GroveExistsRequest{
			NodeID: "root",
		})
		assert.NoError(t, err)
		assert.True(t, existsRes.Exists)

		// Create child1
		child1Position := bullet_interface.ChildPosition(1.0)
		err = c.GroveCreateNode(bullet_interface.GroveCreateNodeRequest{
			NodeID:   "child1",
			Parent:   (*bullet_interface.NodeID)(stringPtr("root")),
			Position: &child1Position,
			Metadata: nil,
		})
		assert.NoError(t, err)

		// Create child2
		child2Position := bullet_interface.ChildPosition(2.0)
		err = c.GroveCreateNode(bullet_interface.GroveCreateNodeRequest{
			NodeID:   "child2",
			Parent:   (*bullet_interface.NodeID)(stringPtr("root")),
			Position: &child2Position,
			Metadata: nil,
		})
		assert.NoError(t, err)

		// Create grandchild1 under child1
		grandchild1Position := bullet_interface.ChildPosition(1.0)
		err = c.GroveCreateNode(bullet_interface.GroveCreateNodeRequest{
			NodeID:   "grandchild1",
			Parent:   (*bullet_interface.NodeID)(stringPtr("child1")),
			Position: &grandchild1Position,
			Metadata: nil,
		})
		assert.NoError(t, err)

		// Get children of root
		childrenRes, err := c.GroveGetChildren(bullet_interface.GroveGetChildrenRequest{
			NodeID:     "root",
			Pagination: nil,
		})
		assert.NoError(t, err)
		assert.Equal(t, 2, len(childrenRes.Children))
		assert.Contains(t, childrenRes.Children, bullet_interface.NodeID("child1"))
		assert.Contains(t, childrenRes.Children, bullet_interface.NodeID("child2"))

		// Get children of child1
		childrenRes, err = c.GroveGetChildren(bullet_interface.GroveGetChildrenRequest{
			NodeID:     "child1",
			Pagination: nil,
		})
		assert.NoError(t, err)
		assert.Equal(t, 1, len(childrenRes.Children))
		assert.Equal(t, bullet_interface.NodeID("grandchild1"), childrenRes.Children[0])

		// Get node info for child1
		nodeInfoRes, err := c.GroveGetNodeInfo(bullet_interface.GroveGetNodeInfoRequest{
			NodeID: "child1",
		})
		assert.NoError(t, err)
		assert.NotNil(t, nodeInfoRes.NodeInfo)
		assert.Equal(t, bullet_interface.NodeID("child1"), nodeInfoRes.NodeInfo.ID)
		assert.Equal(t, bullet_interface.NodeID("root"), *nodeInfoRes.NodeInfo.Parent)

		// Get ancestors of grandchild1 (should be child1, root)
		ancestorsRes, err := c.GroveGetAncestors(bullet_interface.GroveGetAncestorsRequest{
			NodeID:     "grandchild1",
			Pagination: nil,
		})
		assert.NoError(t, err)
		assert.Equal(t, 2, len(ancestorsRes.Ancestors))
		// Ancestors should be [child1, root] or [root, child1] depending on order
		assert.Contains(t, ancestorsRes.Ancestors, bullet_interface.NodeID("child1"))
		assert.Contains(t, ancestorsRes.Ancestors, bullet_interface.NodeID("root"))

		// Get descendants of root
		descendantsRes, err := c.GroveGetDescendants(bullet_interface.GroveGetDescendantsRequest{
			NodeID:  "root",
			Options: nil,
		})
		assert.NoError(t, err)
		assert.Equal(t, 3, len(descendantsRes.Descendants))

		// Test aggregates
		deltas := bullet_interface.AggregateDeltas{
			"count": 5,
		}
		err = c.GroveApplyAggregateMutation(bullet_interface.GroveApplyAggregateMutationRequest{
			MutationID: "mutation1",
			NodeID:     "child1",
			Deltas:     deltas,
		})
		assert.NoError(t, err)

		// Get local aggregates for child1
		localAggRes, err := c.GroveGetNodeLocalAggregates(bullet_interface.GroveGetNodeLocalAggregatesRequest{
			NodeID: "child1",
		})
		assert.NoError(t, err)
		assert.Equal(t, bullet_interface.AggregateValue(5), localAggRes.Aggregates["count"])

		// Get aggregates with descendants for child1 (should include grandchild1)
		withDescAggRes, err := c.GroveGetNodeWithDescendantsAggregates(bullet_interface.GroveGetNodeWithDescendantsAggregatesRequest{
			NodeID: "child1",
		})
		assert.NoError(t, err)
		assert.Equal(t, bullet_interface.AggregateValue(5), withDescAggRes.Aggregates["count"])

		// Test move node - move grandchild1 to be under child2
		newPosition := bullet_interface.ChildPosition(1.0)
		err = c.GroveMoveNode(bullet_interface.GroveMoveNodeRequest{
			NodeID:      "grandchild1",
			NewParent:   (*bullet_interface.NodeID)(stringPtr("child2")),
			NewPosition: &newPosition,
		})
		assert.NoError(t, err)

		// Verify grandchild1 is now under child2
		childrenRes, err = c.GroveGetChildren(bullet_interface.GroveGetChildrenRequest{
			NodeID:     "child2",
			Pagination: nil,
		})
		assert.NoError(t, err)
		assert.Equal(t, 1, len(childrenRes.Children))
		assert.Equal(t, bullet_interface.NodeID("grandchild1"), childrenRes.Children[0])

		// Delete grandchild1 (soft delete)
		err = c.GroveDeleteNode(bullet_interface.GroveDeleteNodeRequest{
			NodeID: "grandchild1",
			Soft:   true,
		})
		assert.NoError(t, err)

		// Verify grandchild1 no longer exists
		existsRes, err = c.GroveExists(bullet_interface.GroveExistsRequest{
			NodeID: "grandchild1",
		})
		assert.NoError(t, err)
		assert.False(t, existsRes.Exists)
	}
}

// Helper function to convert string to *string
func stringPtr(s string) *string {
	return &s
}
