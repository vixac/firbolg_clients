package bullet_interface

// VX:TODO This is missing many depot bullet calls
type DepotClientInterface interface {
	DepotInsertOne(req DepotRequest) error
	DepotGetMany(req DepotGetManyRequest) (*DepotGetManyResponse, error)
	DepotUpsertMany(req []DepotRequest) error
	DepotDeleteOne(req DepotDeleteRequest) error
}

type WayFinderClientInterface interface {
	WayFinderInsertOne(req WayFinderPutRequest) (int64, error)
	WayFinderQueryByPrefix(req WayFinderPrefixQueryRequest) ([]WayFinderQueryItem, error)
	WayFinderGetOne(req WayFinderGetOneRequest) (*WayFinderItem, error)
}

type TrackClientInterface interface {
	TrackGetMany(req TrackGetManyRequest) (*TrackGetManyResponse, error)
	TrackInsertOne(bucketID int32, key string, value int64, tag *int64, metric *float64) error
	TrackDeleteMany(req TrackDeleteMany) error
	TrackGetManyByPrefix(req TrackGetItemsByPrefixRequest) (*TrackGetManyResponse, error)
	TrackGetByManyPrefixes(req TrackGetItemsbyManyPrefixesRequest) (*TrackGetManyResponse, error)
}

type GroveClientInterface interface {
	GroveCreateNode(req GroveCreateNodeRequest) error
	GroveDeleteNode(req GroveDeleteNodeRequest) error
	GroveMoveNode(req GroveMoveNodeRequest) error
	GroveExists(req GroveExistsRequest) (*GroveExistsResponse, error)
	GroveGetNodeInfo(req GroveGetNodeInfoRequest) (*GroveGetNodeInfoResponse, error)
	GroveGetChildren(req GroveGetChildrenRequest) (*GroveGetChildrenResponse, error)
	GroveGetAncestors(req GroveGetAncestorsRequest) (*GroveGetAncestorsResponse, error)
	GroveGetDescendants(req GroveGetDescendantsRequest) (*GroveGetDescendantsResponse, error)
	GroveApplyAggregateMutation(req GroveApplyAggregateMutationRequest) error
	GroveGetNodeLocalAggregates(req GroveGetNodeLocalAggregatesRequest) (*GroveGetAggregatesResponse, error)
	GroveGetNodeWithDescendantsAggregates(req GroveGetNodeWithDescendantsAggregatesRequest) (*GroveGetAggregatesResponse, error)
}

type BulletClientInterface interface {
	TrackClientInterface
	DepotClientInterface
	WayFinderClientInterface
	GroveClientInterface
}
