package local_bullet

import (
	store "github.com/vixac/bullet/store/store_interface"
	"github.com/vixac/firbolg_clients/bullet/bullet_interface"
)

func (l *LocalBullet) GroveCreateNode(req bullet_interface.GroveCreateNodeRequest) error {
	return l.Store.CreateNode(
		l.Space,
		"VX:TODO TREE ID",
		store.NodeID(req.NodeID),
		(*store.NodeID)(req.Parent),
		(*store.ChildPosition)(req.Position),
		(*store.NodeMetadata)(req.Metadata),
	)
}

func (l *LocalBullet) GroveDeleteNode(req bullet_interface.GroveDeleteNodeRequest) error {
	return l.Store.DeleteNode(l.Space, "VX:TODO TREE ID", store.NodeID(req.NodeID), req.Soft)
}

func (l *LocalBullet) GroveMoveNode(req bullet_interface.GroveMoveNodeRequest) error {
	return l.Store.MoveNode(
		l.Space,
		"VX:TODO TREE ID",
		store.NodeID(req.NodeID),
		(*store.NodeID)(req.NewParent),
		(*store.ChildPosition)(req.NewPosition),
	)
}

func (l *LocalBullet) GroveExists(req bullet_interface.GroveExistsRequest) (*bullet_interface.GroveExistsResponse, error) {
	exists, err := l.Store.Exists(l.Space, "VX:TODO TREE ID", store.NodeID(req.NodeID))
	if err != nil {
		return nil, err
	}
	return &bullet_interface.GroveExistsResponse{
		Exists: exists,
	}, nil
}

func (l *LocalBullet) GroveGetNodeInfo(req bullet_interface.GroveGetNodeInfoRequest) (*bullet_interface.GroveGetNodeInfoResponse, error) {
	nodeInfo, err := l.Store.GetNodeInfo(l.Space, "VX:TODO TREE ID", store.NodeID(req.NodeID))
	if err != nil {
		return nil, err
	}
	if nodeInfo == nil {
		return &bullet_interface.GroveGetNodeInfoResponse{
			NodeInfo: nil,
		}, nil
	}
	return &bullet_interface.GroveGetNodeInfoResponse{
		NodeInfo: &bullet_interface.NodeInfo{
			ID:       bullet_interface.NodeID(nodeInfo.ID),
			Parent:   (*bullet_interface.NodeID)(nodeInfo.Parent),
			Position: (*bullet_interface.ChildPosition)(nodeInfo.Position),
			Depth:    nodeInfo.Depth,
			Metadata: (*bullet_interface.NodeMetadata)(nodeInfo.Metadata),
		},
	}, nil
}

func (l *LocalBullet) GroveGetChildren(req bullet_interface.GroveGetChildrenRequest) (*bullet_interface.GroveGetChildrenResponse, error) {
	var pagination *store.PaginationParams
	if req.Pagination != nil {
		pagination = &store.PaginationParams{
			Limit:  req.Pagination.Limit,
			Cursor: req.Pagination.Cursor,
		}
	}
	children, paginationResult, err := l.Store.GetChildren(l.Space, "VX:TODO TREE ID", store.NodeID(req.NodeID), pagination)
	if err != nil {
		return nil, err
	}
	// Convert []store.NodeID to []bullet_interface.NodeID
	childrenConverted := make([]bullet_interface.NodeID, len(children))
	for i, child := range children {
		childrenConverted[i] = bullet_interface.NodeID(child)
	}
	var paginationResultConverted *bullet_interface.PaginationResult
	if paginationResult != nil {
		paginationResultConverted = &bullet_interface.PaginationResult{
			NextCursor: paginationResult.NextCursor,
		}
	}
	return &bullet_interface.GroveGetChildrenResponse{
		Children:   childrenConverted,
		Pagination: paginationResultConverted,
	}, nil
}

func (l *LocalBullet) GroveGetAncestors(req bullet_interface.GroveGetAncestorsRequest) (*bullet_interface.GroveGetAncestorsResponse, error) {
	var pagination *store.PaginationParams
	if req.Pagination != nil {
		pagination = &store.PaginationParams{
			Limit:  req.Pagination.Limit,
			Cursor: req.Pagination.Cursor,
		}
	}
	ancestors, paginationResult, err := l.Store.GetAncestors(l.Space, "VX:TODO TREE ID", store.NodeID(req.NodeID), pagination)
	if err != nil {
		return nil, err
	}
	// Convert []store.NodeID to []bullet_interface.NodeID
	ancestorsConverted := make([]bullet_interface.NodeID, len(ancestors))
	for i, ancestor := range ancestors {
		ancestorsConverted[i] = bullet_interface.NodeID(ancestor)
	}
	var paginationResultConverted *bullet_interface.PaginationResult
	if paginationResult != nil {
		paginationResultConverted = &bullet_interface.PaginationResult{
			NextCursor: paginationResult.NextCursor,
		}
	}
	return &bullet_interface.GroveGetAncestorsResponse{
		Ancestors:  ancestorsConverted,
		Pagination: paginationResultConverted,
	}, nil
}

func (l *LocalBullet) GroveGetDescendants(req bullet_interface.GroveGetDescendantsRequest) (*bullet_interface.GroveGetDescendantsResponse, error) {
	var options *store.DescendantOptions
	if req.Options != nil {
		var pagination *store.PaginationParams
		if req.Options.Pagination != nil {
			pagination = &store.PaginationParams{
				Limit:  req.Options.Pagination.Limit,
				Cursor: req.Options.Pagination.Cursor,
			}
		}
		options = &store.DescendantOptions{
			MaxDepth:     req.Options.MaxDepth,
			IncludeDepth: req.Options.IncludeDepth,
			BreadthFirst: req.Options.BreadthFirst,
			Pagination:   pagination,
		}
	}
	descendants, paginationResult, err := l.Store.GetDescendants(l.Space, "VX:TODO TREE ID", store.NodeID(req.NodeID), options)
	if err != nil {
		return nil, err
	}
	// Convert []store.NodeWithDepth to []bullet_interface.NodeWithDepth
	descendantsConverted := make([]bullet_interface.NodeWithDepth, len(descendants))
	for i, descendant := range descendants {
		descendantsConverted[i] = bullet_interface.NodeWithDepth{
			NodeID: bullet_interface.NodeID(descendant.NodeID),
			Depth:  descendant.Depth,
		}
	}
	var paginationResultConverted *bullet_interface.PaginationResult
	if paginationResult != nil {
		paginationResultConverted = &bullet_interface.PaginationResult{
			NextCursor: paginationResult.NextCursor,
		}
	}
	return &bullet_interface.GroveGetDescendantsResponse{
		Descendants: descendantsConverted,
		Pagination:  paginationResultConverted,
	}, nil
}

func (l *LocalBullet) GroveApplyAggregateMutation(req bullet_interface.GroveApplyAggregateMutationRequest) error {
	// Convert bullet_interface.AggregateDeltas to store.AggregateDeltas
	deltas := make(store.AggregateDeltas)
	for k, v := range req.Deltas {
		deltas[store.AggregateKey(k)] = store.AggregateValue(v)
	}
	return l.Store.ApplyAggregateMutation(
		l.Space,
		"VX:TODO TREE ID",
		store.MutationID(req.MutationID),
		store.NodeID(req.NodeID),
		deltas,
	)
}

func (l *LocalBullet) GroveGetNodeLocalAggregates(req bullet_interface.GroveGetNodeLocalAggregatesRequest) (*bullet_interface.GroveGetAggregatesResponse, error) {
	aggregates, err := l.Store.GetNodeLocalAggregates(l.Space, "VX:TODO TREE ID", store.NodeID(req.NodeID))
	if err != nil {
		return nil, err
	}
	// Convert map[store.AggregateKey]store.AggregateValue to map[bullet_interface.AggregateKey]bullet_interface.AggregateValue
	aggregatesConverted := make(map[bullet_interface.AggregateKey]bullet_interface.AggregateValue)
	for k, v := range aggregates {
		aggregatesConverted[bullet_interface.AggregateKey(k)] = bullet_interface.AggregateValue(v)
	}
	return &bullet_interface.GroveGetAggregatesResponse{
		Aggregates: aggregatesConverted,
	}, nil
}

func (l *LocalBullet) GroveGetNodeWithDescendantsAggregates(req bullet_interface.GroveGetNodeWithDescendantsAggregatesRequest) (*bullet_interface.GroveGetAggregatesResponse, error) {
	aggregates, err := l.Store.GetNodeWithDescendantsAggregates(l.Space, "VX:TODO TREE ID", store.NodeID(req.NodeID))
	if err != nil {
		return nil, err
	}
	// Convert map[store.AggregateKey]store.AggregateValue to map[bullet_interface.AggregateKey]bullet_interface.AggregateValue
	aggregatesConverted := make(map[bullet_interface.AggregateKey]bullet_interface.AggregateValue)
	for k, v := range aggregates {
		aggregatesConverted[bullet_interface.AggregateKey(k)] = bullet_interface.AggregateValue(v)
	}
	return &bullet_interface.GroveGetAggregatesResponse{
		Aggregates: aggregatesConverted,
	}, nil
}
