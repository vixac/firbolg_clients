package rest_bullet

import (
	"encoding/json"
	"fmt"

	bullet_interface "github.com/vixac/firbolg_clients/bullet/bullet_interface"
	util "github.com/vixac/firbolg_clients/bullet/util"
)

type WayFinderClient struct {
	Client *util.FirbolgClient
}

func NewWayFinderClient(baseURL string, appId int64) *WayFinderClient {
	firbolg := util.NewFirbolgClient(baseURL, appId)
	return &WayFinderClient{
		Client: firbolg,
	}
}

func (c *WayFinderClient) WayFinderInsertOne(req bullet_interface.WayFinderPutRequest) (int64, error) {
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

// WayFinderGetOne sends a GET request to retrieve a WayFinder item by bucket/key
func (c *WayFinderClient) WayFinderGetOne(req bullet_interface.WayFinderGetOneRequest) (*bullet_interface.WayFinderItem, error) {
	bodyBytes, err := json.Marshal(req)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal request: %w", err)
	}

	if c.Client == nil {
		return nil, fmt.Errorf("FirbolgClient is nil")
	}
	resp, err := c.Client.PostReq("/wayfinder/get-one", bodyBytes)
	if err != nil {
		return nil, fmt.Errorf("request failed: %w", err)
	}

	var result struct {
		Item bullet_interface.WayFinderItem `json:"item"`
	}
	if err := json.Unmarshal(resp, &result); err != nil {
		return nil, fmt.Errorf("failed to unmarshal response: %w, message body was '%s'", err, string(resp))
	}
	return &result.Item, nil
}

func (c *WayFinderClient) WayFinderQueryByPrefix(req bullet_interface.WayFinderPrefixQueryRequest) ([]bullet_interface.WayFinderQueryItem, error) {
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
		Items []bullet_interface.WayFinderQueryItem `json:"items"`
	}
	if err := json.Unmarshal(resp, &result); err != nil {
		return nil, fmt.Errorf("failed to unmarshal response: %w, message body was '%s'", err, string(resp))
	}
	return result.Items, nil
}
