package bullet

import (
	"encoding/json"
	"fmt"

	util "github.com/vixac/firbolg_clients/util"
)

type WayFinderClientInterface interface {
	WayFinderInsertOne(req WayFinderPutRequest) (int64, error)
	WayFinderQueryByPrefix(req WayFinderPrefixQueryRequest) ([]WayFinderQueryItem, error)
}

type WayFinderClient struct {
	Client *util.FirbolgClient
}

func (c *WayFinderClient) WayFinderInsertOne(req WayFinderPutRequest) (int64, error) {
	bodyBytes, err := json.Marshal(req)
	if err != nil {
		return 0, fmt.Errorf("failed to marshal request: %w", err)
	}

	if c.Client == nil {
		return 0, fmt.Errorf("FirbolgClient is nil")
	}
	resp, err := c.Client.PostReq("/wayfinder/insert-one", bodyBytes)
	if err != nil {
		return 0, fmt.Errorf("request failed: %w", err)
	}

	var result struct {
		ItemId int64 `json:"itemId"`
	}
	if err := json.Unmarshal(resp, &result); err != nil {
		return 0, fmt.Errorf("failed to unmarshal response: %w, message body was '%s'", err, string(resp))
	}
	return result.ItemId, nil
}

func (c *WayFinderClient) WayFinderQueryByPrefix(req WayFinderPrefixQueryRequest) ([]WayFinderQueryItem, error) {
	bodyBytes, err := json.Marshal(req)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal request: %w", err)
	}

	if c.Client == nil {
		return nil, fmt.Errorf("FirbolgClient is nil")
	}
	resp, err := c.Client.PostReq("/wayfinder/query-by-prefix", bodyBytes)
	if err != nil {
		return nil, fmt.Errorf("request failed: %w", err)
	}
	var result struct {
		Items []WayFinderQueryItem `json:"items"`
	}
	if err := json.Unmarshal(resp, &result); err != nil {
		return nil, fmt.Errorf("failed to unmarshal response: %w, message body was '%s'", err, string(resp))
	}
	return result.Items, nil
}
