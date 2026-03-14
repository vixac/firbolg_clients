package bullet_interface

// Types imported from store_interface
type NodeID string
type TreeID string
type AggregateKey string
type MutationID string
type AggregateValue int64
type AggregateDeltas map[AggregateKey]AggregateValue
type ChildPosition float64
type NodeMetadata map[string]interface{}

// Pagination
type PaginationParams struct {
	Limit  int     `json:"limit"`
	Cursor *string `json:"cursor,omitempty"`
}

type PaginationResult struct {
	NextCursor *string `json:"nextCursor,omitempty"`
}

// Node structures
type NodeInfo struct {
	ID       NodeID         `json:"id"`
	Parent   *NodeID        `json:"parent,omitempty"`
	Position *ChildPosition `json:"position,omitempty"`
	Depth    int            `json:"depth"`
	Metadata *NodeMetadata  `json:"metadata,omitempty"`
}

type NodeWithDepth struct {
	NodeID NodeID `json:"nodeId"`
	Depth  int    `json:"depth"`
}

// Query options
type DescendantOptions struct {
	MaxDepth     *int              `json:"maxDepth,omitempty"`
	IncludeDepth bool              `json:"includeDepth"`
	BreadthFirst bool              `json:"breadthFirst"`
	Pagination   *PaginationParams `json:"pagination,omitempty"`
}

// Request structures
type GroveCreateNodeRequest struct {
	NodeID   NodeID         `json:"nodeId"`
	TreeID   TreeID         `json:"treeId"`
	Parent   *NodeID        `json:"parent,omitempty"`
	Position *ChildPosition `json:"position,omitempty"`
	Metadata *NodeMetadata  `json:"metadata,omitempty"`
}

type GroveDeleteNodeRequest struct {
	NodeID NodeID `json:"nodeId"`
	TreeID TreeID `json:"treeId"`
	Soft   bool   `json:"soft"`
}

type GroveMoveNodeRequest struct {
	NodeID      NodeID         `json:"nodeId"`
	TreeID      TreeID         `json:"treeId"`
	NewParent   *NodeID        `json:"newParent,omitempty"`
	NewPosition *ChildPosition `json:"newPosition,omitempty"`
}

type GroveExistsRequest struct {
	NodeID NodeID `json:"nodeId"`
	TreeID TreeID `json:"treeId"`
}

type GroveGetNodeInfoRequest struct {
	NodeID NodeID `json:"nodeId"`
	TreeID TreeID `json:"treeId"`
}

type GroveGetChildrenRequest struct {
	NodeID     NodeID            `json:"nodeId"`
	TreeID     TreeID            `json:"treeId"`
	Pagination *PaginationParams `json:"pagination,omitempty"`
}

type GroveGetAncestorsRequest struct {
	NodeID     NodeID            `json:"nodeId"`
	TreeID     TreeID            `json:"treeId"`
	Pagination *PaginationParams `json:"pagination,omitempty"`
}

type GroveGetDescendantsRequest struct {
	NodeID  NodeID             `json:"nodeId"`
	TreeID  TreeID             `json:"treeId"`
	Options *DescendantOptions `json:"options,omitempty"`
}

type GroveApplyAggregateMutationRequest struct {
	MutationID MutationID      `json:"mutationId"`
	NodeID     NodeID          `json:"nodeId"`
	TreeID     TreeID          `json:"treeId"`
	Deltas     AggregateDeltas `json:"deltas"`
}

type GroveGetNodeLocalAggregatesRequest struct {
	NodeID NodeID `json:"nodeId"`
	TreeID TreeID `json:"treeId"`
}

type GroveGetNodeWithDescendantsAggregatesRequest struct {
	NodeID NodeID `json:"nodeId"`
	TreeID TreeID `json:"treeId"`
}

// Response structures
type GroveExistsResponse struct {
	Exists bool `json:"exists"`
}

type GroveGetNodeInfoResponse struct {
	NodeInfo *NodeInfo `json:"nodeInfo,omitempty"`
}

type GroveGetChildrenResponse struct {
	Children   []NodeID          `json:"children"`
	Pagination *PaginationResult `json:"pagination,omitempty"`
}

type GroveGetAncestorsResponse struct {
	Ancestors  []NodeID          `json:"ancestors"`
	Pagination *PaginationResult `json:"pagination,omitempty"`
}

type GroveGetDescendantsResponse struct {
	Descendants []NodeWithDepth   `json:"descendants"`
	Pagination  *PaginationResult `json:"pagination,omitempty"`
}

type GroveGetAggregatesResponse struct {
	Aggregates map[AggregateKey]AggregateValue `json:"aggregates"`
}
