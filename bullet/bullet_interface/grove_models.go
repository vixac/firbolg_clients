package bullet_interface

// Types imported from store_interface
type NodeID string
type AggregateKey string
type MutationID string
type AggregateValue int64
type AggregateDeltas map[AggregateKey]AggregateValue
type ChildPosition float64
type NodeMetadata map[string]interface{}

// Pagination
type PaginationParams struct {
	Limit  int
	Cursor *string
}

type PaginationResult struct {
	NextCursor *string
}

// Node structures
type NodeInfo struct {
	ID       NodeID
	Parent   *NodeID
	Position *ChildPosition
	Depth    int
	Metadata *NodeMetadata
}

type NodeWithDepth struct {
	NodeID NodeID
	Depth  int
}

// Query options
type DescendantOptions struct {
	MaxDepth     *int
	IncludeDepth bool
	BreadthFirst bool
	Pagination   *PaginationParams
}

// Request structures
type GroveCreateNodeRequest struct {
	NodeID   NodeID
	Parent   *NodeID
	Position *ChildPosition
	Metadata *NodeMetadata
}

type GroveDeleteNodeRequest struct {
	NodeID NodeID
	Soft   bool
}

type GroveMoveNodeRequest struct {
	NodeID      NodeID
	NewParent   *NodeID
	NewPosition *ChildPosition
}

type GroveExistsRequest struct {
	NodeID NodeID
}

type GroveGetNodeInfoRequest struct {
	NodeID NodeID
}

type GroveGetChildrenRequest struct {
	NodeID     NodeID
	Pagination *PaginationParams
}

type GroveGetAncestorsRequest struct {
	NodeID     NodeID
	Pagination *PaginationParams
}

type GroveGetDescendantsRequest struct {
	NodeID  NodeID
	Options *DescendantOptions
}

type GroveApplyAggregateMutationRequest struct {
	MutationID MutationID
	NodeID     NodeID
	Deltas     AggregateDeltas
}

type GroveGetNodeLocalAggregatesRequest struct {
	NodeID NodeID
}

type GroveGetNodeWithDescendantsAggregatesRequest struct {
	NodeID NodeID
}

// Response structures
type GroveExistsResponse struct {
	Exists bool
}

type GroveGetNodeInfoResponse struct {
	NodeInfo *NodeInfo
}

type GroveGetChildrenResponse struct {
	Children   []NodeID
	Pagination *PaginationResult
}

type GroveGetAncestorsResponse struct {
	Ancestors  []NodeID
	Pagination *PaginationResult
}

type GroveGetDescendantsResponse struct {
	Descendants []NodeWithDepth
	Pagination  *PaginationResult
}

type GroveGetAggregatesResponse struct {
	Aggregates map[AggregateKey]AggregateValue
}
