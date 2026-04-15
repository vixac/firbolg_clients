package test_suite

import (
	"fmt"
	"net/http/httptest"
	"path/filepath"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/vixac/bullet/api"
	ram "github.com/vixac/bullet/store/ram"
	sqlite_store "github.com/vixac/bullet/store/sqlite"
	"github.com/vixac/bullet/store/store_interface"
	"github.com/vixac/firbolg_clients/bullet/bullet_interface"
	local_bullet "github.com/vixac/firbolg_clients/bullet/local_bullet"
	"github.com/vixac/firbolg_clients/bullet/rest_bullet"
)

type clientPair struct {
	name  string
	local bullet_interface.BulletClientInterface
	rest  bullet_interface.BulletClientInterface
}

func buildClientPairs(t *testing.T) []clientPair {
	t.Helper()

	space := store_interface.TenancySpace{
		AppId:     12,
		TenancyId: 100,
	}

	buildPair := func(name string, store store_interface.Store) clientPair {
		server := newBulletServer(t, store)
		t.Cleanup(server.Close)

		return clientPair{
			name: name,
			local: &local_bullet.LocalBullet{
				Store: store,
				Space: space,
			},
			rest: rest_bullet.NewRestClient(server.URL, space, rest_bullet.WithHTTPClient(server.Client())),
		}
	}

	sqlitePath := filepath.Join(t.TempDir(), "test-sqlite.db")
	sqliteStore, err := sqlite_store.NewSQLiteStore(sqlitePath)
	require.NoError(t, err)

	return []clientPair{
		buildPair("ram", ram.NewRamStore()),
		buildPair("sqlite", sqliteStore),
	}
}

func newBulletServer(t *testing.T, store store_interface.Store) *httptest.Server {
	t.Helper()

	gin.SetMode(gin.TestMode)
	engine := gin.New()
	api.SetupTrackRouter(store, "/track", engine)
	api.SetupDepotRouter(store, "/depot", engine)
	api.SetupGroveRouter(store, "/grove", engine)
	return httptest.NewServer(engine)
}

func TestTrack(t *testing.T) {
	for _, pair := range buildClientPairs(t) {
		t.Run(pair.name, func(t *testing.T) {
			err := pair.local.TrackInsertOne(1, "testKey", 1234, nil, nil)
			require.NoError(t, err)

			tag := int64(7)
			metric := 3.5
			err = pair.rest.TrackInsertOne(1, "testKey_2", 12345, &tag, &metric)
			require.NoError(t, err)

			err = pair.local.TrackInsertOne(1, "not_a_testKey3", 123456, nil, nil)
			require.NoError(t, err)

			getReq := bullet_interface.TrackGetManyRequest{
				Buckets: []bullet_interface.TrackGetKeys{
					{BucketID: 1, Keys: []string{"testKey", "testKey_2", "missing"}},
				},
			}

			localMany, err := pair.local.TrackGetMany(getReq)
			require.NoError(t, err)
			restMany, err := pair.rest.TrackGetMany(getReq)
			require.NoError(t, err)
			assert.Equal(t, localMany, restMany)
			assert.Equal(t, int64(1234), restMany.Values[1]["testKey"].Value)
			assert.Equal(t, int64(12345), restMany.Values[1]["testKey_2"].Value)
			assert.Equal(t, []string{"missing"}, restMany.Missing["1"])

			prefixReq := bullet_interface.TrackGetItemsbyManyPrefixesRequest{
				BucketID: 1,
				Prefixes: []string{"testKey_", "not_"},
			}
			localPrefix, err := pair.local.TrackGetByManyPrefixes(prefixReq)
			require.NoError(t, err)
			restPrefix, err := pair.rest.TrackGetByManyPrefixes(prefixReq)
			require.NoError(t, err)
			assert.Equal(t, localPrefix, restPrefix)
			assert.Equal(t, int64(12345), restPrefix.Values[1]["testKey_2"].Value)
			assert.Equal(t, int64(123456), restPrefix.Values[1]["not_a_testKey3"].Value)

			deleteReq := bullet_interface.TrackDeleteMany{
				Values: []bullet_interface.TrackDeleteValue{
					{BucketID: 1, Key: "testKey"},
				},
			}
			err = pair.rest.TrackDeleteMany(deleteReq)
			require.NoError(t, err)

			afterDeleteReq := bullet_interface.TrackGetManyRequest{
				Buckets: []bullet_interface.TrackGetKeys{
					{BucketID: 1, Keys: []string{"testKey", "testKey_2"}},
				},
			}
			localAfterDelete, err := pair.local.TrackGetMany(afterDeleteReq)
			require.NoError(t, err)
			restAfterDelete, err := pair.rest.TrackGetMany(afterDeleteReq)
			require.NoError(t, err)
			assert.Equal(t, localAfterDelete, restAfterDelete)
			assert.Contains(t, restAfterDelete.Missing["1"], "testKey")
		})
	}
}

func TestDepot(t *testing.T) {
	for _, pair := range buildClientPairs(t) {
		t.Run(pair.name, func(t *testing.T) {
			const bucket = int32(1)

			createResp, err := pair.local.DepotCreate(bullet_interface.DepotCreateRequest{
				BucketID: bucket,
				Value:    "value1",
			})
			require.NoError(t, err)
			id1 := createResp.ID

			getResp, err := pair.rest.DepotGetOne(bullet_interface.DepotGetRequest{ID: id1})
			require.NoError(t, err)
			assert.Equal(t, "value1", getResp.Value)

			createManyResp, err := pair.rest.DepotCreateMany(bullet_interface.DepotCreateManyRequest{
				BucketID: bucket,
				Values:   []string{"value2", "value3"},
			})
			require.NoError(t, err)
			require.Len(t, createManyResp.IDs, 2)
			id2, id3 := createManyResp.IDs[0], createManyResp.IDs[1]

			err = pair.local.DepotUpdate(bullet_interface.DepotUpdateRequest{ID: id1, Value: "updated_value1"})
			require.NoError(t, err)

			localMany, err := pair.local.DepotGetMany(bullet_interface.DepotGetManyRequest{IDs: []int64{id1, id3, -999}})
			require.NoError(t, err)
			restMany, err := pair.rest.DepotGetMany(bullet_interface.DepotGetManyRequest{IDs: []int64{id1, id3, -999}})
			require.NoError(t, err)
			assert.Equal(t, localMany, restMany)
			assert.Equal(t, "updated_value1", restMany.Values[id1])
			assert.Equal(t, "value3", restMany.Values[id3])

			localAll, err := pair.local.DepotGetAllByBucket(bullet_interface.DepotBucketRequest{BucketID: bucket})
			require.NoError(t, err)
			restAll, err := pair.rest.DepotGetAllByBucket(bullet_interface.DepotBucketRequest{BucketID: bucket})
			require.NoError(t, err)
			assert.Equal(t, localAll, restAll)
			assert.Len(t, restAll.Values, 3)

			err = pair.rest.DepotDelete(bullet_interface.DepotDeleteRequest{ID: id2})
			require.NoError(t, err)

			localAfterDelete, err := pair.local.DepotGetAllByBucket(bullet_interface.DepotBucketRequest{BucketID: bucket})
			require.NoError(t, err)
			restAfterDelete, err := pair.rest.DepotGetAllByBucket(bullet_interface.DepotBucketRequest{BucketID: bucket})
			require.NoError(t, err)
			assert.Equal(t, localAfterDelete, restAfterDelete)
			assert.Len(t, restAfterDelete.Values, 2)

			err = pair.local.DepotDeleteByBucket(bullet_interface.DepotBucketRequest{BucketID: bucket})
			require.NoError(t, err)

			localAfterBucketDelete, err := pair.local.DepotGetAllByBucket(bullet_interface.DepotBucketRequest{BucketID: bucket})
			require.NoError(t, err)
			restAfterBucketDelete, err := pair.rest.DepotGetAllByBucket(bullet_interface.DepotBucketRequest{BucketID: bucket})
			require.NoError(t, err)
			assert.Equal(t, localAfterBucketDelete, restAfterBucketDelete)
			assert.Len(t, restAfterBucketDelete.Values, 0)
		})
	}
}

func TestGrove(t *testing.T) {
	for _, pair := range buildClientPairs(t) {
		t.Run(pair.name, func(t *testing.T) {
			treeID := bullet_interface.TreeID("test-tree")

			err := pair.local.GroveCreateNode(bullet_interface.GroveCreateNodeRequest{
				TreeID: treeID,
				NodeID: "root",
			})
			require.NoError(t, err)

			existsRes, err := pair.rest.GroveExists(bullet_interface.GroveExistsRequest{
				TreeID: treeID,
				NodeID: "root",
			})
			require.NoError(t, err)
			assert.True(t, existsRes.Exists)

			child1Position := bullet_interface.ChildPosition(1.0)
			err = pair.rest.GroveCreateNode(bullet_interface.GroveCreateNodeRequest{
				TreeID:   treeID,
				NodeID:   "child1",
				Parent:   nodeIDPtr("root"),
				Position: &child1Position,
			})
			require.NoError(t, err)

			child2Position := bullet_interface.ChildPosition(2.0)
			err = pair.local.GroveCreateNode(bullet_interface.GroveCreateNodeRequest{
				TreeID:   treeID,
				NodeID:   "child2",
				Parent:   nodeIDPtr("root"),
				Position: &child2Position,
			})
			require.NoError(t, err)

			grandchildPosition := bullet_interface.ChildPosition(1.0)
			err = pair.rest.GroveCreateNode(bullet_interface.GroveCreateNodeRequest{
				TreeID:   treeID,
				NodeID:   "grandchild1",
				Parent:   nodeIDPtr("child1"),
				Position: &grandchildPosition,
			})
			require.NoError(t, err)

			localChildren, err := pair.local.GroveGetChildren(bullet_interface.GroveGetChildrenRequest{
				TreeID: treeID,
				NodeID: "root",
			})
			require.NoError(t, err)
			restChildren, err := pair.rest.GroveGetChildren(bullet_interface.GroveGetChildrenRequest{
				TreeID: treeID,
				NodeID: "root",
			})
			require.NoError(t, err)
			assert.Equal(t, localChildren, restChildren)
			assert.Contains(t, restChildren.Children, bullet_interface.NodeID("child1"))
			assert.Contains(t, restChildren.Children, bullet_interface.NodeID("child2"))

			localInfo, err := pair.local.GroveGetNodeInfo(bullet_interface.GroveGetNodeInfoRequest{
				TreeID: treeID,
				NodeID: "child1",
			})
			require.NoError(t, err)
			restInfo, err := pair.rest.GroveGetNodeInfo(bullet_interface.GroveGetNodeInfoRequest{
				TreeID: treeID,
				NodeID: "child1",
			})
			require.NoError(t, err)
			assert.Equal(t, localInfo, restInfo)

			localAncestors, err := pair.local.GroveGetAncestors(bullet_interface.GroveGetAncestorsRequest{
				TreeID: treeID,
				NodeID: "grandchild1",
			})
			require.NoError(t, err)
			restAncestors, err := pair.rest.GroveGetAncestors(bullet_interface.GroveGetAncestorsRequest{
				TreeID: treeID,
				NodeID: "grandchild1",
			})
			require.NoError(t, err)
			assert.ElementsMatch(t, localAncestors.Ancestors, restAncestors.Ancestors)
			assert.Equal(t, localAncestors.Pagination, restAncestors.Pagination)

			localDescendants, err := pair.local.GroveGetDescendants(bullet_interface.GroveGetDescendantsRequest{
				TreeID: treeID,
				NodeID: "root",
			})
			require.NoError(t, err)
			restDescendants, err := pair.rest.GroveGetDescendants(bullet_interface.GroveGetDescendantsRequest{
				TreeID: treeID,
				NodeID: "root",
			})
			require.NoError(t, err)
			assert.ElementsMatch(t, localDescendants.Descendants, restDescendants.Descendants)
			assert.Equal(t, localDescendants.Pagination, restDescendants.Pagination)
			assert.Len(t, restDescendants.Descendants, 3)

			err = pair.local.GroveApplyAggregateMutation(bullet_interface.GroveApplyAggregateMutationRequest{
				TreeID:     treeID,
				MutationID: "mutation1",
				NodeID:     "child1",
				Deltas: bullet_interface.AggregateDeltas{
					"count": 5,
				},
			})
			require.NoError(t, err)

			localLocalAgg, err := pair.local.GroveGetNodeLocalAggregates(bullet_interface.GroveGetNodeLocalAggregatesRequest{
				TreeID: treeID,
				NodeID: "child1",
			})
			require.NoError(t, err)
			restLocalAgg, err := pair.rest.GroveGetNodeLocalAggregates(bullet_interface.GroveGetNodeLocalAggregatesRequest{
				TreeID: treeID,
				NodeID: "child1",
			})
			require.NoError(t, err)
			assert.Equal(t, localLocalAgg, restLocalAgg)

			localSubtreeAgg, err := pair.local.GroveGetNodeWithDescendantsAggregates(bullet_interface.GroveGetNodeWithDescendantsAggregatesRequest{
				TreeID: treeID,
				NodeID: "child1",
			})
			require.NoError(t, err)
			restSubtreeAgg, err := pair.rest.GroveGetNodeWithDescendantsAggregates(bullet_interface.GroveGetNodeWithDescendantsAggregatesRequest{
				TreeID: treeID,
				NodeID: "child1",
			})
			require.NoError(t, err)
			assert.Equal(t, localSubtreeAgg, restSubtreeAgg)

			err = pair.rest.GroveMoveNode(bullet_interface.GroveMoveNodeRequest{
				TreeID:      treeID,
				NodeID:      "grandchild1",
				NewParent:   nodeIDPtr("child2"),
				NewPosition: &grandchildPosition,
			})
			require.NoError(t, err)

			localChild2Children, err := pair.local.GroveGetChildren(bullet_interface.GroveGetChildrenRequest{
				TreeID: treeID,
				NodeID: "child2",
			})
			require.NoError(t, err)
			restChild2Children, err := pair.rest.GroveGetChildren(bullet_interface.GroveGetChildrenRequest{
				TreeID: treeID,
				NodeID: "child2",
			})
			require.NoError(t, err)
			assert.Equal(t, localChild2Children, restChild2Children)
			assert.Equal(t, []bullet_interface.NodeID{"grandchild1"}, restChild2Children.Children)

			bulkNodes := []bullet_interface.NodeID{"child1", "child2", "missing"}
			localBulkAncestors, err := pair.local.GroveGetAncestorsBulk(bullet_interface.GroveGetAncestorsBulkRequest{
				TreeID:  treeID,
				NodeIDs: bulkNodes,
			})
			require.NoError(t, err)
			restBulkAncestors, err := pair.rest.GroveGetAncestorsBulk(bullet_interface.GroveGetAncestorsBulkRequest{
				TreeID:  treeID,
				NodeIDs: bulkNodes,
			})
			require.NoError(t, err)
			assert.Equal(t, localBulkAncestors, restBulkAncestors)

			localBulkLocalAgg, err := pair.local.GroveGetNodeLocalAggregatesBulk(bullet_interface.GroveGetNodeLocalAggregatesBulkRequest{
				TreeID:  treeID,
				NodeIDs: bulkNodes,
			})
			require.NoError(t, err)
			restBulkLocalAgg, err := pair.rest.GroveGetNodeLocalAggregatesBulk(bullet_interface.GroveGetNodeLocalAggregatesBulkRequest{
				TreeID:  treeID,
				NodeIDs: bulkNodes,
			})
			require.NoError(t, err)
			assert.Equal(t, localBulkLocalAgg, restBulkLocalAgg)

			localBulkSubtreeAgg, err := pair.local.GroveGetNodeWithDescendantsAggregatesBulk(bullet_interface.GroveGetNodeWithDescendantsAggregatesBulkRequest{
				TreeID:  treeID,
				NodeIDs: bulkNodes,
			})
			require.NoError(t, err)
			restBulkSubtreeAgg, err := pair.rest.GroveGetNodeWithDescendantsAggregatesBulk(bullet_interface.GroveGetNodeWithDescendantsAggregatesBulkRequest{
				TreeID:  treeID,
				NodeIDs: bulkNodes,
			})
			require.NoError(t, err)
			assert.Equal(t, localBulkSubtreeAgg, restBulkSubtreeAgg)

			err = pair.local.GroveDeleteNode(bullet_interface.GroveDeleteNodeRequest{
				TreeID: treeID,
				NodeID: "grandchild1",
				Soft:   true,
			})
			require.NoError(t, err)

			existsAfterDeleteLocal, err := pair.local.GroveExists(bullet_interface.GroveExistsRequest{
				TreeID: treeID,
				NodeID: "grandchild1",
			})
			require.NoError(t, err)
			existsAfterDeleteRest, err := pair.rest.GroveExists(bullet_interface.GroveExistsRequest{
				TreeID: treeID,
				NodeID: "grandchild1",
			})
			require.NoError(t, err)
			assert.Equal(t, existsAfterDeleteLocal, existsAfterDeleteRest)
			assert.False(t, existsAfterDeleteRest.Exists)
		})
	}
}

func nodeIDPtr(value string) *bullet_interface.NodeID {
	nodeID := bullet_interface.NodeID(value)
	return &nodeID
}

func Example_buildClientPairs() {
	fmt.Println("isomorphic client pairs")
}
