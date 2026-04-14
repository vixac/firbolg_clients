package rest_bullet

import (
	"encoding/json"
	"fmt"
	"net/http"

	bullet_model "github.com/vixac/bullet/model"
	"github.com/vixac/firbolg_clients/bullet/bullet_interface"
	"github.com/vixac/firbolg_clients/bullet/util"
)

func (c *RestClient) GroveCreateNode(req bullet_interface.GroveCreateNodeRequest) error {
	bodyBytes, err := util.MarshalJSONBody(bullet_model.GroveCreateNodeRequest{
		NodeID:   string(req.NodeID),
		ParentID: nodeIDPtrToString(req.Parent),
		Position: childPositionToFloat(req.Position),
		Metadata: nodeMetadataToMap(req.Metadata),
	})
	if err != nil {
		return err
	}

	_, err = c.PostReq(groveTreePath(req.TreeID)+"/nodes", bodyBytes, http.StatusCreated)
	return err
}

func (c *RestClient) GroveDeleteNode(req bullet_interface.GroveDeleteNodeRequest) error {
	path := fmt.Sprintf("%s/nodes/%s", groveTreePath(req.TreeID), req.NodeID)
	if req.Soft {
		path += "?soft=true"
	}
	_, err := c.DeleteReq(path, nil, http.StatusNoContent)
	return err
}

func (c *RestClient) GroveMoveNode(req bullet_interface.GroveMoveNodeRequest) error {
	bodyBytes, err := util.MarshalJSONBody(bullet_model.GroveMoveNodeRequest{
		NewParentID: nodeIDPtrToString(req.NewParent),
		NewPosition: childPositionToFloat(req.NewPosition),
	})
	if err != nil {
		return err
	}

	_, err = c.PatchReq(fmt.Sprintf("%s/nodes/%s", groveTreePath(req.TreeID), req.NodeID), bodyBytes, http.StatusOK)
	return err
}

func (c *RestClient) GroveExists(req bullet_interface.GroveExistsRequest) (*bullet_interface.GroveExistsResponse, error) {
	resp, err := c.GetReq(fmt.Sprintf("%s/nodes/%s/exists", groveTreePath(req.TreeID), req.NodeID), http.StatusOK)
	if err != nil {
		return nil, err
	}

	var result bullet_model.GroveExistsResponse
	if err := json.Unmarshal(resp, &result); err != nil {
		return nil, fmt.Errorf("failed to unmarshal response: %w", err)
	}
	return &bullet_interface.GroveExistsResponse{Exists: result.Exists}, nil
}

func (c *RestClient) GroveGetNodeInfo(req bullet_interface.GroveGetNodeInfoRequest) (*bullet_interface.GroveGetNodeInfoResponse, error) {
	resp, err := c.GetReq(fmt.Sprintf("%s/nodes/%s", groveTreePath(req.TreeID), req.NodeID), http.StatusOK)
	if err != nil {
		return nil, err
	}

	var result bullet_model.GroveNodeInfoResponse
	if err := json.Unmarshal(resp, &result); err != nil {
		return nil, fmt.Errorf("failed to unmarshal response: %w", err)
	}

	return &bullet_interface.GroveGetNodeInfoResponse{
		NodeInfo: &bullet_interface.NodeInfo{
			ID:       bullet_interface.NodeID(result.ID),
			Parent:   stringPtrToNodeID(result.ParentID),
			Position: floatPtrToChildPosition(result.Position),
			Depth:    result.Depth,
			Metadata: mapToNodeMetadata(result.Metadata),
		},
	}, nil
}

func (c *RestClient) GroveGetChildren(req bullet_interface.GroveGetChildrenRequest) (*bullet_interface.GroveGetChildrenResponse, error) {
	resp, err := c.GetReq(fmt.Sprintf("%s/nodes/%s/children", groveTreePath(req.TreeID), req.NodeID), http.StatusOK)
	if err != nil {
		return nil, err
	}

	var result bullet_model.GroveChildrenResponse
	if err := json.Unmarshal(resp, &result); err != nil {
		return nil, fmt.Errorf("failed to unmarshal response: %w", err)
	}
	return &bullet_interface.GroveGetChildrenResponse{
		Children:   stringsToNodeIDs(result.Children),
		Pagination: &bullet_interface.PaginationResult{},
	}, nil
}

func (c *RestClient) GroveGetAncestors(req bullet_interface.GroveGetAncestorsRequest) (*bullet_interface.GroveGetAncestorsResponse, error) {
	resp, err := c.GetReq(fmt.Sprintf("%s/nodes/%s/ancestors", groveTreePath(req.TreeID), req.NodeID), http.StatusOK)
	if err != nil {
		return nil, err
	}

	var result bullet_model.GroveAncestorsResponse
	if err := json.Unmarshal(resp, &result); err != nil {
		return nil, fmt.Errorf("failed to unmarshal response: %w", err)
	}
	return &bullet_interface.GroveGetAncestorsResponse{
		Ancestors:  stringsToNodeIDs(result.Ancestors),
		Pagination: &bullet_interface.PaginationResult{},
	}, nil
}

func (c *RestClient) GroveGetDescendants(req bullet_interface.GroveGetDescendantsRequest) (*bullet_interface.GroveGetDescendantsResponse, error) {
	resp, err := c.GetReq(fmt.Sprintf("%s/nodes/%s/descendants", groveTreePath(req.TreeID), req.NodeID), http.StatusOK)
	if err != nil {
		return nil, err
	}

	var result bullet_model.GroveDescendantsResponse
	if err := json.Unmarshal(resp, &result); err != nil {
		return nil, fmt.Errorf("failed to unmarshal response: %w", err)
	}

	descendants := make([]bullet_interface.NodeWithDepth, len(result.Descendants))
	for i, descendant := range result.Descendants {
		descendants[i] = bullet_interface.NodeWithDepth{
			NodeID: bullet_interface.NodeID(descendant.NodeID),
			Depth:  descendant.Depth,
		}
	}

	return &bullet_interface.GroveGetDescendantsResponse{
		Descendants: descendants,
		Pagination:  &bullet_interface.PaginationResult{},
	}, nil
}

func (c *RestClient) GroveApplyAggregateMutation(req bullet_interface.GroveApplyAggregateMutationRequest) error {
	bodyBytes, err := util.MarshalJSONBody(bullet_model.GroveApplyMutationRequest{
		MutationID: string(req.MutationID),
		Deltas:     aggregateDeltasToMap(req.Deltas),
	})
	if err != nil {
		return err
	}

	_, err = c.PostReq(fmt.Sprintf("%s/nodes/%s/mutations", groveTreePath(req.TreeID), req.NodeID), bodyBytes, http.StatusOK)
	return err
}

func (c *RestClient) GroveGetNodeLocalAggregates(req bullet_interface.GroveGetNodeLocalAggregatesRequest) (*bullet_interface.GroveGetAggregatesResponse, error) {
	resp, err := c.GetReq(fmt.Sprintf("%s/nodes/%s/aggregates/local", groveTreePath(req.TreeID), req.NodeID), http.StatusOK)
	if err != nil {
		return nil, err
	}
	return parseAggregateResponse(resp)
}

func (c *RestClient) GroveGetNodeWithDescendantsAggregates(req bullet_interface.GroveGetNodeWithDescendantsAggregatesRequest) (*bullet_interface.GroveGetAggregatesResponse, error) {
	resp, err := c.GetReq(fmt.Sprintf("%s/nodes/%s/aggregates", groveTreePath(req.TreeID), req.NodeID), http.StatusOK)
	if err != nil {
		return nil, err
	}
	return parseAggregateResponse(resp)
}

func (c *RestClient) GroveGetAncestorsBulk(req bullet_interface.GroveGetAncestorsBulkRequest) (*bullet_interface.GroveGetAncestorsBulkResponse, error) {
	bodyBytes, err := util.MarshalJSONBody(bullet_model.GroveBulkNodesRequest{
		NodeIDs: nodeIDsToStrings(req.NodeIDs),
	})
	if err != nil {
		return nil, err
	}
	resp, err := c.PostReq(groveTreePath(req.TreeID)+"/bulk/ancestors", bodyBytes, http.StatusOK)
	if err != nil {
		return nil, err
	}

	var result bullet_model.GroveAncestorsBulkResponse
	if err := json.Unmarshal(resp, &result); err != nil {
		return nil, fmt.Errorf("failed to unmarshal response: %w", err)
	}

	ancestors := make(map[bullet_interface.NodeID][]bullet_interface.NodeID, len(result.Ancestors))
	for nodeID, items := range result.Ancestors {
		ancestors[bullet_interface.NodeID(nodeID)] = stringsToNodeIDs(items)
	}

	return &bullet_interface.GroveGetAncestorsBulkResponse{
		Ancestors:    ancestors,
		MissingNodes: stringsToNodeIDs(result.Missing),
	}, nil
}

func (c *RestClient) GroveGetNodeLocalAggregatesBulk(req bullet_interface.GroveGetNodeLocalAggregatesBulkRequest) (*bullet_interface.GroveGetNodeLocalAggregatesBulkResponse, error) {
	resp, err := c.groveAggregatesBulk(req.TreeID, "/bulk/aggregates/local", req.NodeIDs)
	if err != nil {
		return nil, err
	}
	return &bullet_interface.GroveGetNodeLocalAggregatesBulkResponse{
		Aggregates:   resp.Aggregates,
		MissingNodes: resp.MissingNodes,
	}, nil
}

func (c *RestClient) GroveGetNodeWithDescendantsAggregatesBulk(req bullet_interface.GroveGetNodeWithDescendantsAggregatesBulkRequest) (*bullet_interface.GroveGetNodeWithDescendantsAggregatesBulkResponse, error) {
	resp, err := c.groveAggregatesBulk(req.TreeID, "/bulk/aggregates", req.NodeIDs)
	if err != nil {
		return nil, err
	}
	return &bullet_interface.GroveGetNodeWithDescendantsAggregatesBulkResponse{
		Aggregates:   resp.Aggregates,
		MissingNodes: resp.MissingNodes,
	}, nil
}

func (c *RestClient) groveAggregatesBulk(treeID bullet_interface.TreeID, suffix string, nodeIDs []bullet_interface.NodeID) (*bullet_interface.GroveGetNodeWithDescendantsAggregatesBulkResponse, error) {
	bodyBytes, err := util.MarshalJSONBody(bullet_model.GroveBulkNodesRequest{
		NodeIDs: nodeIDsToStrings(nodeIDs),
	})
	if err != nil {
		return nil, err
	}
	resp, err := c.PostReq(groveTreePath(treeID)+suffix, bodyBytes, http.StatusOK)
	if err != nil {
		return nil, err
	}

	var result bullet_model.GroveAggregatesBulkResponse
	if err := json.Unmarshal(resp, &result); err != nil {
		return nil, fmt.Errorf("failed to unmarshal response: %w", err)
	}

	return &bullet_interface.GroveGetNodeWithDescendantsAggregatesBulkResponse{
		Aggregates:   aggregateBulkMapsToInterface(result.Aggregates),
		MissingNodes: stringsToNodeIDs(result.Missing),
	}, nil
}

func parseAggregateResponse(resp []byte) (*bullet_interface.GroveGetAggregatesResponse, error) {
	var result bullet_model.GroveAggregatesResponse
	if err := json.Unmarshal(resp, &result); err != nil {
		return nil, fmt.Errorf("failed to unmarshal response: %w", err)
	}
	return &bullet_interface.GroveGetAggregatesResponse{
		Aggregates: aggregateMapToInterface(result.Aggregates),
	}, nil
}

func groveTreePath(treeID bullet_interface.TreeID) string {
	return "/grove/trees/" + string(treeID)
}

func nodeIDPtrToString(nodeID *bullet_interface.NodeID) *string {
	if nodeID == nil {
		return nil
	}
	value := string(*nodeID)
	return &value
}

func stringPtrToNodeID(value *string) *bullet_interface.NodeID {
	if value == nil {
		return nil
	}
	nodeID := bullet_interface.NodeID(*value)
	return &nodeID
}

func childPositionToFloat(position *bullet_interface.ChildPosition) *float64 {
	if position == nil {
		return nil
	}
	value := float64(*position)
	return &value
}

func floatPtrToChildPosition(value *float64) *bullet_interface.ChildPosition {
	if value == nil {
		return nil
	}
	position := bullet_interface.ChildPosition(*value)
	return &position
}

func nodeMetadataToMap(metadata *bullet_interface.NodeMetadata) map[string]interface{} {
	if metadata == nil {
		return nil
	}
	return map[string]interface{}(*metadata)
}

func mapToNodeMetadata(metadata map[string]interface{}) *bullet_interface.NodeMetadata {
	if metadata == nil {
		return nil
	}
	value := bullet_interface.NodeMetadata(metadata)
	return &value
}

func aggregateDeltasToMap(deltas bullet_interface.AggregateDeltas) map[string]int64 {
	result := make(map[string]int64, len(deltas))
	for key, value := range deltas {
		result[string(key)] = int64(value)
	}
	return result
}

func aggregateMapToInterface(values map[string]int64) map[bullet_interface.AggregateKey]bullet_interface.AggregateValue {
	result := make(map[bullet_interface.AggregateKey]bullet_interface.AggregateValue, len(values))
	for key, value := range values {
		result[bullet_interface.AggregateKey(key)] = bullet_interface.AggregateValue(value)
	}
	return result
}

func aggregateBulkMapsToInterface(values map[string]map[string]int64) map[bullet_interface.NodeID]map[bullet_interface.AggregateKey]bullet_interface.AggregateValue {
	result := make(map[bullet_interface.NodeID]map[bullet_interface.AggregateKey]bullet_interface.AggregateValue, len(values))
	for nodeID, aggregates := range values {
		result[bullet_interface.NodeID(nodeID)] = aggregateMapToInterface(aggregates)
	}
	return result
}

func stringsToNodeIDs(values []string) []bullet_interface.NodeID {
	result := make([]bullet_interface.NodeID, len(values))
	for i, value := range values {
		result[i] = bullet_interface.NodeID(value)
	}
	return result
}

func nodeIDsToStrings(values []bullet_interface.NodeID) []string {
	result := make([]string, len(values))
	for i, value := range values {
		result[i] = string(value)
	}
	return result
}
